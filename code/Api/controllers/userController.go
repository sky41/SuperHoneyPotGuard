package controllers

import (
	"net/http"
	"superhoneypotguard/database"
	"superhoneypotguard/middleware"
	"superhoneypotguard/models"
	"superhoneypotguard/utils"
	"time"

	"github.com/gin-gonic/gin"
)

type UserController struct{}

func NewUserController() *UserController {
	return &UserController{}
}

func (ctrl *UserController) GetList(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("pageSize", "10")
	username := c.Query("username")
	status := c.Query("status")

	query := database.DB.Model(&models.User{})

	if username != "" {
		query = query.Where("username LIKE ?", "%"+username+"%")
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	offset := (parseInt(page) - 1) * parseInt(pageSize)

	var users []models.User
	query.Offset(offset).Limit(parseInt(pageSize)).Order("created_at DESC").Find(&users)

	for i := range users {
		var roles []models.Role
		database.DB.Raw(`
			SELECT r.id, r.role_name, r.role_code
			FROM roles r
			INNER JOIN user_roles ur ON r.id = ur.role_id
			WHERE ur.user_id = ?
		`, users[i].ID).Scan(&roles)
		users[i].Roles = roles
	}

	utils.SuccessResponse(c, models.PaginatedResponse{
		List:     users,
		Total:    total,
		Page:     parseInt(page),
		PageSize: parseInt(pageSize),
	})
}

func (ctrl *UserController) GetById(c *gin.Context) {
	id := c.Param("id")

	var user models.User
	if err := database.DB.Where("id = ?", id).First(&user).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "用户不存在")
		return
	}

	var roles []models.Role
	database.DB.Raw(`
		SELECT r.id, r.role_name, r.role_code
		FROM roles r
		INNER JOIN user_roles ur ON r.id = ur.role_id
		WHERE ur.user_id = ?
	`, user.ID).Scan(&roles)
	user.Roles = roles

	utils.SuccessResponse(c, user)
}

func (ctrl *UserController) Create(c *gin.Context) {
	var req models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "参数验证失败")
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

	currentUser := middleware.GetCurrentUser(c)
	status := 1
	if req.Status != nil {
		status = *req.Status
	}

	user := models.User{
		Username:  req.Username,
		Password:  hashedPassword,
		Email:     req.Email,
		Phone:     req.Phone,
		RealName:  req.RealName,
		Status:    status,
		CreatedBy: &currentUser.UserID,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "创建用户失败")
		return
	}

	if len(req.RoleIDs) > 0 {
		for _, roleId := range req.RoleIDs {
			database.DB.Create(&models.UserRole{
				UserID:    user.ID,
				RoleID:    roleId,
				CreatedBy: &currentUser.UserID,
			})
		}
	}

	utils.SuccessResponse(c, gin.H{
		"id":       user.ID,
		"username": user.Username,
	})
}

func (ctrl *UserController) Update(c *gin.Context) {
	id := c.Param("id")

	var user models.User
	if err := database.DB.Where("id = ?", id).First(&user).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "用户不存在")
		return
	}

	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "参数验证失败")
		return
	}

	currentUser := middleware.GetCurrentUser(c)
	updates := map[string]interface{}{
		"email":     req.Email,
		"phone":     req.Phone,
		"real_name": req.RealName,
		"updated_by": currentUser.UserID,
	}

	if req.Status != nil {
		updates["status"] = *req.Status
	}

	if err := database.DB.Model(&user).Updates(updates).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "更新用户失败")
		return
	}

	if req.RoleIDs != nil {
		database.DB.Where("user_id = ?", user.ID).Delete(&models.UserRole{})

		for _, roleId := range *req.RoleIDs {
			database.DB.Create(&models.UserRole{
				UserID:    user.ID,
				RoleID:    roleId,
				CreatedBy: &currentUser.UserID,
			})
		}
	}

	utils.SuccessResponse(c, nil)
}

func (ctrl *UserController) Delete(c *gin.Context) {
	id := c.Param("id")
	currentUser := middleware.GetCurrentUser(c)

	if parseInt(id) == currentUser.UserID {
		utils.ErrorResponse(c, http.StatusBadRequest, "不能删除当前登录用户")
		return
	}

	var user models.User
	if err := database.DB.Where("id = ?", id).First(&user).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "用户不存在")
		return
	}

	if err := database.DB.Delete(&user).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "删除用户失败")
		return
	}

	utils.SuccessResponse(c, nil)
}

func (ctrl *UserController) UpdateStatus(c *gin.Context) {
	id := c.Param("id")
	currentUser := middleware.GetCurrentUser(c)

	if parseInt(id) == currentUser.UserID {
		utils.ErrorResponse(c, http.StatusBadRequest, "不能修改当前登录用户状态")
		return
	}

	var req models.UpdateUserStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "参数验证失败")
		return
	}

	var user models.User
	if err := database.DB.Where("id = ?", id).First(&user).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "用户不存在")
		return
	}

	if err := database.DB.Model(&user).Updates(map[string]interface{}{
		"status":     req.Status,
		"updated_by": currentUser.UserID,
	}).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "更新用户状态失败")
		return
	}

	utils.SuccessResponse(c, nil)
}

func (ctrl *UserController) ResetPassword(c *gin.Context) {
	id := c.Param("id")
	currentUser := middleware.GetCurrentUser(c)

	var req models.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "参数验证失败")
		return
	}

	var user models.User
	if err := database.DB.Where("id = ?", id).First(&user).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "用户不存在")
		return
	}

	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "密码加密失败")
		return
	}

	if err := database.DB.Model(&user).Updates(map[string]interface{}{
		"password":   hashedPassword,
		"updated_by": currentUser.UserID,
	}).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "重置密码失败")
		return
	}

	utils.SuccessResponse(c, nil)
}

func parseInt(s string) int {
	var result int
	for _, c := range s {
		if c >= '0' && c <= '9' {
			result = result*10 + int(c-'0')
		}
	}
	return result
}
