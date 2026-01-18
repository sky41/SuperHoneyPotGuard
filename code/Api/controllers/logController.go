package controllers

import (
	"superhoneypotguard/database"
	"superhoneypotguard/models"
	"superhoneypotguard/utils"

	"github.com/gin-gonic/gin"
)

type LogController struct{}

func NewLogController() *LogController {
	return &LogController{}
}

func (ctrl *LogController) GetList(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("pageSize", "10")
	username := c.Query("username")
	operation := c.Query("operation")
	status := c.Query("status")

	query := database.DB.Model(&models.OperationLog{})

	if username != "" {
		query = query.Where("username LIKE ?", "%"+username+"%")
	}

	if operation != "" {
		query = query.Where("operation LIKE ?", "%"+operation+"%")
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	offset := (parseInt(page) - 1) * parseInt(pageSize)

	var logs []models.OperationLog
	query.Offset(offset).Limit(parseInt(pageSize)).Order("created_at DESC").Find(&logs)

	utils.SuccessResponse(c, models.PaginatedResponse{
		List:     logs,
		Total:    total,
		Page:     parseInt(page),
		PageSize: parseInt(pageSize),
	})
}

func (ctrl *LogController) GetById(c *gin.Context) {
	id := c.Param("id")

	var log models.OperationLog
	if err := database.DB.Where("id = ?", id).First(&log).Error; err != nil {
		utils.ErrorResponse(c, 404, "日志不存在")
		return
	}

	utils.SuccessResponse(c, log)
}

func (ctrl *LogController) Delete(c *gin.Context) {
	id := c.Param("id")

	var log models.OperationLog
	if err := database.DB.Where("id = ?", id).First(&log).Error; err != nil {
		utils.ErrorResponse(c, 404, "日志不存在")
		return
	}

	if err := database.DB.Delete(&log).Error; err != nil {
		utils.ErrorResponse(c, 500, "删除日志失败")
		return
	}

	utils.SuccessResponse(c, nil)
}

func (ctrl *LogController) Clear(c *gin.Context) {
	database.DB.Exec("DELETE FROM operation_logs")

	utils.SuccessResponse(c, nil)
}
