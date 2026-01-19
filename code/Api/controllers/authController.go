package controllers

import (
	"net/http"
	"superhoneypotguard/database"
	"superhoneypotguard/middleware"
	"superhoneypotguard/models"
	"superhoneypotguard/services"
	"superhoneypotguard/utils"
	"time"

	"github.com/gin-gonic/gin"
)

type AuthController struct{}

func NewAuthController() *AuthController {
	return &AuthController{}
}

func (ctrl *AuthController) SendVerificationCode(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "参数验证失败")
		return
	}

	emailService := services.NewEmailService(database.DB)
	if err := emailService.SendVerificationCode(req.Email); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "发送验证码失败: "+err.Error())
		return
	}

	utils.SuccessResponse(c, nil)
}

func (ctrl *AuthController) Register(c *gin.Context) {
	var req struct {
		Username string  `json:"username" binding:"required,min=3,max=50"`
		Password string  `json:"password" binding:"required,min=6"`
		Email    string  `json:"email" binding:"required,email"`
		Code     string  `json:"code" binding:"required,len=6"`
		Phone    *string `json:"phone"`
		RealName *string `json:"realName"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "参数验证失败")
		return
	}

	if req.Email == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "邮箱不能为空")
		return
	}

	emailService := services.NewEmailService(database.DB)
	if !emailService.VerifyCode(req.Email, req.Code) {
		utils.ErrorResponse(c, http.StatusBadRequest, "验证码错误或已过期")
		return
	}

	var count int64
	database.DB.Model(&models.User{}).Where("username = ? OR email = ?", req.Username, req.Email).Count(&count)
	if count > 0 {
		utils.ErrorResponse(c, http.StatusBadRequest, "用户名或邮箱已存在")
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "密码加密失败")
		return
	}

	user := models.User{
		Username: req.Username,
		Password: hashedPassword,
		Email:    &req.Email,
		Phone:    req.Phone,
		RealName: req.RealName,
		Status:   1,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "注册失败")
		return
	}

	var defaultRole models.Role
	database.DB.Where("role_code = ?", "USER").First(&defaultRole)
	if defaultRole.ID > 0 {
		database.DB.Create(&models.UserRole{
			UserID: user.ID,
			RoleID: defaultRole.ID,
		})
	}

	utils.SuccessResponse(c, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
	})
}

func (ctrl *AuthController) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "参数验证失败")
		return
	}

	var user models.User
	if err := database.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "用户名或密码错误")
		return
	}

	if user.Status != 1 {
		utils.ErrorResponse(c, http.StatusForbidden, "账号已被禁用")
		return
	}

	if !utils.ComparePassword(req.Password, user.Password) {
		utils.ErrorResponse(c, http.StatusUnauthorized, "用户名或密码错误")
		return
	}

	now := time.Now()
	ip := utils.GetClientIP(c)
	database.DB.Model(&user).Updates(map[string]interface{}{
		"last_login_time": &now,
		"last_login_ip":   &ip,
	})

	var roles []models.Role
	database.DB.Raw(`
		SELECT r.id, r.role_name, r.role_code
		FROM roles r
		INNER JOIN user_roles ur ON r.id = ur.role_id
		WHERE ur.user_id = ? AND r.status = 1
	`, user.ID).Scan(&roles)

	roleCodes := make([]string, 0, len(roles))
	for _, role := range roles {
		roleCodes = append(roleCodes, role.RoleCode)
	}

	var permissions []struct {
		PermissionCode string `json:"permissionCode"`
	}
	database.DB.Raw(`
		SELECT DISTINCT p.permission_code
		FROM permissions p
		INNER JOIN role_permissions rp ON p.id = rp.permission_id
		INNER JOIN user_roles ur ON rp.role_id = ur.role_id
		WHERE ur.user_id = ? AND p.status = 1
	`, user.ID).Scan(&permissions)

	permissionCodes := make([]string, 0, len(permissions))
	for _, perm := range permissions {
		permissionCodes = append(permissionCodes, perm.PermissionCode)
	}

	claims := &models.Claims{
		UserID:      user.ID,
		Username:    user.Username,
		Roles:       roleCodes,
		Permissions: permissionCodes,
	}

	token, err := utils.GenerateToken(claims)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "生成令牌失败")
		return
	}

	utils.SuccessResponse(c, gin.H{
		"token": token,
		"user": gin.H{
			"id":          user.ID,
			"username":    user.Username,
			"email":       user.Email,
			"realName":    user.RealName,
			"roles":       roles,
			"permissions": permissions,
		},
	})
}

func (ctrl *AuthController) Logout(c *gin.Context) {
	utils.SuccessResponse(c, nil)
}

func (ctrl *AuthController) GetCurrentUser(c *gin.Context) {
	user := middleware.GetCurrentUser(c)

	var dbUser models.User
	if err := database.DB.Where("id = ?", user.UserID).First(&dbUser).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "用户不存在")
		return
	}

	var roles []models.Role
	database.DB.Raw(`
		SELECT r.id, r.role_name, r.role_code
		FROM roles r
		INNER JOIN user_roles ur ON r.id = ur.role_id
		WHERE ur.user_id = ? AND r.status = 1
	`, user.UserID).Scan(&roles)

	var permissions []models.Permission
	database.DB.Raw(`
		SELECT DISTINCT p.id, p.permission_name, p.permission_code, p.permission_type
		FROM permissions p
		INNER JOIN role_permissions rp ON p.id = rp.permission_id
		INNER JOIN user_roles ur ON rp.role_id = ur.role_id
		WHERE ur.user_id = ? AND p.status = 1
	`, user.UserID).Scan(&permissions)

	utils.SuccessResponse(c, gin.H{
		"user": gin.H{
			"id":            dbUser.ID,
			"username":      dbUser.Username,
			"email":         dbUser.Email,
			"realName":      dbUser.RealName,
			"status":        dbUser.Status,
			"lastLoginTime": dbUser.LastLoginTime,
			"createdAt":     dbUser.CreatedAt,
			"updatedAt":     dbUser.UpdatedAt,
		},
		"roles":       roles,
		"permissions": permissions,
	})
}
