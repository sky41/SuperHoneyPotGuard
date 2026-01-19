package controllers

import (
	"net/http"
	"superhoneypotguard/database"
	"superhoneypotguard/models"
	"superhoneypotguard/services"
	"superhoneypotguard/utils"

	"github.com/gin-gonic/gin"
)

type PasswordController struct{}

func NewPasswordController() *PasswordController {
	return &PasswordController{}
}

func (ctrl *PasswordController) SendResetPasswordCode(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "参数验证失败")
		return
	}

	emailService := services.NewEmailService(database.DB)
	if err := emailService.SendResetPasswordCode(req.Email); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "发送验证码失败: "+err.Error())
		return
	}

	utils.SuccessResponse(c, nil)
}

func (ctrl *PasswordController) ResetPassword(c *gin.Context) {
	var req struct {
		Email       string `json:"email" binding:"required,email"`
		Code        string `json:"code" binding:"required,len=6"`
		NewPassword string `json:"newPassword" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "参数验证失败")
		return
	}

	emailService := services.NewEmailService(database.DB)
	if !emailService.VerifyCode(req.Email, req.Code) {
		utils.ErrorResponse(c, http.StatusBadRequest, "验证码错误或已过期")
		return
	}

	var user models.User
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "用户不存在")
		return
	}

	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "密码加密失败")
		return
	}

	user.Password = hashedPassword
	if err := database.DB.Save(&user).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "密码重置失败")
		return
	}

	utils.SuccessResponse(c, nil)
}
