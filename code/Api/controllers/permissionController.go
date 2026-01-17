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

type PermissionController struct{}

func NewPermissionController() *PermissionController {
	return &PermissionController{}
}

func (ctrl *PermissionController) GetTree(c *gin.Context) {
	var permissions []models.Permission
	database.DB.Order("sort_order ASC, id ASC").Find(&permissions)

	tree := buildTree(permissions, 0)

	utils.SuccessResponse(c, tree)
}

func buildTree(permissions []models.Permission, parentId int) []models.Permission {
	var result []models.Permission
	for _, p := range permissions {
		if p.ParentID == parentId {
			p.Children = buildTree(permissions, p.ID)
			result = append(result, p)
		}
	}
	return result
}

func (ctrl *PermissionController) GetAll(c *gin.Context) {
	var permissions []models.Permission
	database.DB.Where("status = 1").Order("sort_order ASC, id ASC").Find(&permissions)

	utils.SuccessResponse(c, permissions)
}

func (ctrl *PermissionController) GetById(c *gin.Context) {
	id := c.Param("id")

	var permission models.Permission
	if err := database.DB.Where("id = ?", id).First(&permission).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "权限不存在")
		return
	}

	utils.SuccessResponse(c, permission)
}

func (ctrl *PermissionController) Create(c *gin.Context) {
	var req models.CreatePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "参数验证失败")
		return
	}

	var count int64
	database.DB.Model(&models.Permission{}).Where("permission_code = ?", req.PermissionCode).Count(&count)
	if count > 0 {
		utils.ErrorResponse(c, http.StatusBadRequest, "权限编码已存在")
		return
	}

	parentId := 0
	if req.ParentID != nil {
		parentId = *req.ParentID
	}

	sortOrder := 0
	if req.SortOrder != nil {
		sortOrder = *req.SortOrder
	}

	status := 1
	if req.Status != nil {
		status = *req.Status
	}

	permission := models.Permission{
		PermissionName: req.PermissionName,
		PermissionCode: req.PermissionCode,
		PermissionType: req.PermissionType,
		ParentID:       parentId,
		Path:           req.Path,
		Component:      req.Component,
		Icon:           req.Icon,
		SortOrder:      sortOrder,
		Description:    req.Description,
		Status:         status,
	}

	if err := database.DB.Create(&permission).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "创建权限失败")
		return
	}

	utils.SuccessResponse(c, gin.H{
		"id":             permission.ID,
		"permissionName": permission.PermissionName,
	})
}

func (ctrl *PermissionController) Update(c *gin.Context) {
	id := c.Param("id")

	var permission models.Permission
	if err := database.DB.Where("id = ?", id).First(&permission).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "权限不存在")
		return
	}

	var req models.UpdatePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "参数验证失败")
		return
	}

	parentId := 0
	if req.ParentID != nil {
		parentId = *req.ParentID
	}

	sortOrder := 0
	if req.SortOrder != nil {
		sortOrder = *req.SortOrder
	}

	status := 1
	if req.Status != nil {
		status = *req.Status
	}

	updates := map[string]interface{}{
		"permission_name": req.PermissionName,
		"permission_type": req.PermissionType,
		"parent_id":       parentId,
		"path":            req.Path,
		"component":       req.Component,
		"icon":            req.Icon,
		"sort_order":      sortOrder,
		"description":     req.Description,
		"status":          status,
	}

	if err := database.DB.Model(&permission).Updates(updates).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "更新权限失败")
		return
	}

	utils.SuccessResponse(c, nil)
}

func (ctrl *PermissionController) Delete(c *gin.Context) {
	id := c.Param("id")

	var count int64
	database.DB.Model(&models.Permission{}).Where("parent_id = ?", id).Count(&count)
	if count > 0 {
		utils.ErrorResponse(c, http.StatusBadRequest, "该权限下还有子权限，无法删除")
		return
	}

	var permission models.Permission
	if err := database.DB.Where("id = ?", id).First(&permission).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "权限不存在")
		return
	}

	if err := database.DB.Delete(&permission).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "删除权限失败")
		return
	}

	utils.SuccessResponse(c, nil)
}

func parseInt(s string) int {
	result, _ := strconv.Atoi(s)
	return result
}
