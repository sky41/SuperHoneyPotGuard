package controllers

import (
	"superhoneypotguard/database"
	"superhoneypotguard/models"
	"superhoneypotguard/utils"

	"github.com/gin-gonic/gin"
)

type DashboardController struct{}

func NewDashboardController() *DashboardController {
	return &DashboardController{}
}

func (ctrl *DashboardController) GetStats(c *gin.Context) {
	var userCount int64
	database.DB.Model(&models.User{}).Count(&userCount)

	var roleCount int64
	database.DB.Model(&models.Role{}).Count(&roleCount)

	var permissionCount int64
	database.DB.Model(&models.Permission{}).Count(&permissionCount)

	var logCount int64
	database.DB.Model(&models.OperationLog{}).Count(&logCount)

	utils.SuccessResponse(c, gin.H{
		"userCount":    userCount,
		"roleCount":    roleCount,
		"permissionCount": permissionCount,
		"logCount":      logCount,
	})
}
