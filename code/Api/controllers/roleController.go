package controllers

import (
	"net/http"
	"superhoneypotguard/database"
	"superhoneypotguard/middleware"
	"superhoneypotguard/models"
	"superhoneypotguard/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RoleController struct{}

func NewRoleController() *RoleController {
	return &RoleController{}
}

func (ctrl *RoleController) GetList(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("pageSize", "10")
	roleName := c.Query("roleName")
	status := c.Query("status")

	query := database.DB.Model(&models.Role{})

	if roleName != "" {
		query = query.Where("role_name LIKE ?", "%"+roleName+"%")
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	offset := (parseInt(page) - 1) * parseInt(pageSize)

	var roles []models.Role
	query.Offset(offset).Limit(parseInt(pageSize)).Order("created_at DESC").Find(&roles)

	for i := range roles {
		var permissions []models.Permission
		database.DB.Raw(`
			SELECT p.id, p.permission_name, p.permission_code, p.permission_type
			FROM permissions p
			INNER JOIN role_permissions rp ON p.id = rp.permission_id
			WHERE rp.role_id = ?
		`, roles[i].ID).Scan(&permissions)
		roles[i].Permissions = permissions
	}

	utils.SuccessResponse(c, models.PaginatedResponse{
		List:     roles,
		Total:    total,
		Page:     parseInt(page),
		PageSize: parseInt(pageSize),
	})
}

func (ctrl *RoleController) GetAll(c *gin.Context) {
	var roles []models.Role
	database.DB.Where("status = 1").Order("id ASC").Find(&roles)

	utils.SuccessResponse(c, roles)
}

func (ctrl *RoleController) GetById(c *gin.Context) {
	id := c.Param("id")

	var role models.Role
	if err := database.DB.Where("id = ?", id).First(&role).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "角色不存在")
		return
	}

	var permissions []models.Permission
	database.DB.Raw(`
		SELECT p.id, p.permission_name, p.permission_code, p.permission_type
		FROM permissions p
		INNER JOIN role_permissions rp ON p.id = rp.permission_id
		WHERE rp.role_id = ?
	`, role.ID).Scan(&permissions)
	role.Permissions = permissions

	utils.SuccessResponse(c, role)
}

func (ctrl *RoleController) Create(c *gin.Context) {
	var req models.CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "参数验证失败")
		return
	}

	var count int64
	database.DB.Model(&models.Role{}).Where("role_name = ? OR role_code = ?", req.RoleName, req.RoleCode).Count(&count)
	if count > 0 {
		utils.ErrorResponse(c, http.StatusBadRequest, "角色名称或角色编码已存在")
		return
	}

	currentUser := middleware.GetCurrentUser(c)
	status := 1
	if req.Status != nil {
		status = *req.Status
	}

	role := models.Role{
		RoleName:    req.RoleName,
		RoleCode:    req.RoleCode,
		Description: req.Description,
		Status:      status,
		CreatedBy:   &currentUser.UserID,
	}

	if err := database.DB.Create(&role).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "创建角色失败")
		return
	}

	if len(req.PermissionIDs) > 0 {
		for _, permissionId := range req.PermissionIDs {
			database.DB.Create(&models.RolePermission{
				RoleID:       role.ID,
				PermissionID: permissionId,
				CreatedBy:    &currentUser.UserID,
			})
		}
	}

	utils.SuccessResponse(c, gin.H{
		"id":       role.ID,
		"roleName": role.RoleName,
	})
}

func (ctrl *RoleController) Update(c *gin.Context) {
	id := c.Param("id")

	var role models.Role
	if err := database.DB.Where("id = ?", id).First(&role).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "角色不存在")
		return
	}

	var req models.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "参数验证失败")
		return
	}

	currentUser := middleware.GetCurrentUser(c)
	updates := map[string]interface{}{
		"role_name":  req.RoleName,
		"description": req.Description,
		"updated_by": currentUser.UserID,
	}

	if req.Status != nil {
		updates["status"] = *req.Status
	}

	if err := database.DB.Model(&role).Updates(updates).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "更新角色失败")
		return
	}

	if req.PermissionIDs != nil {
		database.DB.Where("role_id = ?", role.ID).Delete(&models.RolePermission{})

		for _, permissionId := range *req.PermissionIDs {
			database.DB.Create(&models.RolePermission{
				RoleID:       role.ID,
				PermissionID: permissionId,
				CreatedBy:    &currentUser.UserID,
			})
		}
	}

	utils.SuccessResponse(c, nil)
}

func (ctrl *RoleController) Delete(c *gin.Context) {
	id := c.Param("id")

	var count int64
	database.DB.Model(&models.UserRole{}).Where("role_id = ?", id).Count(&count)
	if count > 0 {
		utils.ErrorResponse(c, http.StatusBadRequest, "该角色下还有用户，无法删除")
		return
	}

	var role models.Role
	if err := database.DB.Where("id = ?", id).First(&role).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "角色不存在")
		return
	}

	if err := database.DB.Delete(&role).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "删除角色失败")
		return
	}

	utils.SuccessResponse(c, nil)
}

func parseInt(s string) int {
	result, _ := strconv.Atoi(s)
	return result
}
