package services

import (
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"strings"
	"time"

	"superhoneypotguard/config"
	"superhoneypotguard/models"

	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
)

type EmailService struct {
	DB *gorm.DB
}

func NewEmailService(db *gorm.DB) *EmailService {
	return &EmailService{DB: db}
}

type VerificationCode struct {
	Code      string
	Email     string
	ExpiresAt time.Time
	SentAt    time.Time
}

var verificationCodes = make(map[string]*VerificationCode)

// 提交目的：优化SMTP连接配置，支持多种端口和加密方式
// 提交内容：添加统一的sendEmail方法，改进TLS/SSL配置，增强日志记录
// 提交时间：2026-01-18

// 提交目的：添加数据库验证码存储，支持频率限制和邮箱检查
// 提交内容：修改EmailService结构，添加数据库支持，实现验证码持久化
// 提交时间：2026-01-19

// InitEmailService 初始化邮件服务，创建验证码表
func InitEmailService(db *gorm.DB) error {
	// 创建验证码表
	if err := db.AutoMigrate(&VerificationCode{}); err != nil {
		log.Printf("创建验证码表失败: %v", err)
		return err
	}
	log.Printf("验证码表初始化成功")
	return nil
}

// sendEmail 统一的邮件发送方法
// isResetPassword: 是否为密码重置验证码（true=重置密码，false=注册）
func (s *EmailService) sendEmail(email, subject, body string, isResetPassword bool) error {
	cfg := config.AppConfig

	log.Printf("准备发送邮件到: %s", email)
	log.Printf("SMTP 配置详情: Host=%s, Port=%s, User=%s", cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUser)

	// 根据场景应用不同的验证码发送频率限制
	if isResetPassword {
		// 密码重置：3分钟内最多10次，超过后需要等待10分钟
		var recentCodes []VerificationCode
		threeMinutesAgo := time.Now().Add(-3 * time.Minute)
		tenMinutesAgo := time.Now().Add(-10 * time.Minute)

		if err := s.DB.Where("email = ? AND sent_at > ? AND sent_at < ?", email, tenMinutesAgo, threeMinutesAgo).Find(&recentCodes).Error; err != nil {
			log.Printf("查询验证码记录失败: %v", err)
			return fmt.Errorf("系统错误，请稍后重试")
		}

		if len(recentCodes) >= 10 {
			// 检查是否有超过10分钟的记录
			var hasOldCode bool
			for _, code := range recentCodes {
				if code.SentAt.Before(threeMinutesAgo) {
					hasOldCode = true
					break
				}
			}

			if hasOldCode {
				// 如果3分钟内有超过10分钟的记录，需要等待
				oldestCode := recentCodes[0]
				waitTime := 10*time.Minute - time.Since(oldestCode.SentAt)
				remainingSeconds := int(waitTime.Seconds())
				log.Printf("密码重置验证码发送过于频繁，请 %d 秒后再试", remainingSeconds)
				return fmt.Errorf("验证码发送过于频繁，请 %d 秒后再试", remainingSeconds)
			} else {
				// 如果3分钟内都是最近10分钟的记录，则已达到限制
				log.Printf("密码重置验证码3分钟内已发送10次，请10分钟后再试")
				return fmt.Errorf("验证码发送过于频繁，请10分钟后再试")
			}
		}
	} else {
		// 注册：检查邮箱是否已注册
		var existingUser models.User
		if err := s.DB.Where("email = ?", email).First(&existingUser).Error; err == nil {
			log.Printf("邮箱 %s 已被注册", email)
			return fmt.Errorf("该邮箱已被注册")
		} else if err != gorm.ErrRecordNotFound {
			log.Printf("查询用户失败: %v", err)
			return fmt.Errorf("系统错误，请稍后重试")
		}

		// 注册：60秒内不能重复发送
		var existingCode VerificationCode
		if err := s.DB.Where("email = ?", email).Order("sent_at DESC").First(&existingCode).Error; err == nil {
			timeSinceLastSent := time.Since(existingCode.SentAt)
			if timeSinceLastSent < 60*time.Second {
				log.Printf("邮箱 %s 最近已发送过验证码，距离上次发送: %v 秒", email, timeSinceLastSent.Seconds())
				return fmt.Errorf("验证码发送过于频繁，请 %d 秒后再试", 60-int(timeSinceLastSent.Seconds()))
			}
		}
	}

	from := cfg.SMTPUser
	to := []string{email}

	port, err := strconv.Atoi(cfg.SMTPPort)
	if err != nil {
		log.Printf("SMTP 端口转换失败: %v", err)
		return fmt.Errorf("SMTP 端口转换失败: %v", err)
	}

	log.Printf("SMTP 端口: %d", port)

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", strings.Join(to, ","))
	m.SetHeader("Subject", subject)
	m.SetHeader("MIME-Version", "1.0")
	m.SetHeader("Content-Type", "text/html; charset=UTF-8")
	m.SetBody("text/html", body)

	d := gomail.NewDialer(cfg.SMTPHost, port, cfg.SMTPUser, cfg.SMTPPassword)

	// 根据端口配置加密方式
	switch port {
	case 465:
		// SSL/TLS 加密
		d.SSL = true
		log.Printf("使用 SSL 加密 (端口 465)")
	case 587:
		// STARTTLS 加密
		d.TLSConfig = &tls.Config{
			InsecureSkipVerify: false,
			ServerName:         cfg.SMTPHost,
		}
		log.Printf("使用 TLS/STARTTLS 加密 (端口 587)")
	case 25:
		// 通常非加密，但一些服务商可能需要 TLS
		d.TLSConfig = &tls.Config{
			InsecureSkipVerify: false,
			ServerName:         cfg.SMTPHost,
		}
		log.Printf("使用 TLS 加密 (端口 25)")
	default:
		// 其他端口默认使用 TLS
		d.TLSConfig = &tls.Config{
			InsecureSkipVerify: false,
			ServerName:         cfg.SMTPHost,
		}
		log.Printf("使用 TLS 加密 (端口 %d)", port)
	}

	log.Printf("开始连接 SMTP 服务器: %s:%d", cfg.SMTPHost, port)

	// 发送邮件
	if err := d.DialAndSend(m); err != nil {
		log.Printf("发送邮件失败: %v", err)
		return fmt.Errorf("发送邮件失败: %v", err)
	}

	log.Printf("邮件发送成功到: %s", email)
	return nil
}

func (s *EmailService) SendVerificationCode(email string) error {
	code := generateVerificationCode()

	subject := "SuperHoneyPotGuard 注册验证码"

	body := fmt.Sprintf(`
		<h2>注册验证码</h2>
		<p>您好，</p>
		<p>您正在注册 SuperHoneyPotGuard 账号。</p>
		<p>您的验证码是：<strong style="font-size: 24px; color: #1890ff;">%s</strong></p>
		<p>验证码有效期为 5 分钟。</p>
		<p>如果这不是您本人操作，请忽略此邮件。</p>
		<p>此邮件由系统自动发送，请勿回复。</p>
	`, code)

	if err := s.sendEmail(email, subject, body, false); err != nil {
		return err
	}

	// 存储验证码到数据库
	verificationCode := &VerificationCode{
		Code:      code,
		Email:     email,
		ExpiresAt: time.Now().Add(5 * time.Minute),
		SentAt:    time.Now(),
	}

	if err := s.DB.Create(verificationCode).Error; err != nil {
		log.Printf("存储验证码失败: %v", err)
		return fmt.Errorf("存储验证码失败: %v", err)
	}

	log.Printf("验证码存储成功: %s -> %s", email, code)
	return nil
}

func (s *EmailService) SendResetPasswordCode(email string) error {
	code := generateVerificationCode()

	subject := "SuperHoneyPotGuard 密码重置验证码"

	body := fmt.Sprintf(`
		<h2>密码重置验证码</h2>
		<p>您好，</p>
		<p>您正在重置 SuperHoneyPotGuard 账号密码。</p>
		<p>您的验证码是：<strong style="font-size: 24px; color: #1890ff;">%s</strong></p>
		<p>验证码有效期为 5 分钟。</p>
		<p>如果这不是您本人操作，请忽略此邮件。</p>
		<p>此邮件由系统自动发送，请勿回复。</p>
	`, code)

	if err := s.sendEmail(email, subject, body, true); err != nil {
		return err
	}

	// 存储验证码到数据库
	verificationCode := &VerificationCode{
		Code:      code,
		Email:     email,
		ExpiresAt: time.Now().Add(5 * time.Minute),
		SentAt:    time.Now(),
	}

	if err := s.DB.Create(verificationCode).Error; err != nil {
		log.Printf("存储验证码失败: %v", err)
		return fmt.Errorf("存储验证码失败: %v", err)
	}

	log.Printf("验证码存储成功: %s -> %s", email, code)
	return nil
}

func (s *EmailService) VerifyCode(email, code string) bool {
	var verificationCode VerificationCode
	if err := s.DB.Where("email = ? AND code = ? AND expires_at > ?", email, code, time.Now()).First(&verificationCode).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Printf("验证码不存在或已过期: %s", email)
		} else {
			log.Printf("查询验证码失败: %v", err)
		}
		return false
	}

	log.Printf("验证码验证成功: %s", email)

	// 验证成功后删除验证码
	if err := s.DB.Delete(&verificationCode).Error; err != nil {
		log.Printf("删除验证码失败: %v", err)
	}

	return true
}

func generateVerificationCode() string {
	code := ""
	max := big.NewInt(10)
	for i := 0; i < 6; i++ {
		n, _ := rand.Int(rand.Reader, max)
		code += fmt.Sprintf("%d", n)
	}
	return code
}

// CleanupExpiredCodes 清理过期的验证码
func (s *EmailService) CleanupExpiredCodes() {
	// 从数据库中删除过期的验证码
	result := s.DB.Where("expires_at < ?", time.Now()).Delete(&VerificationCode{})
	if result.Error != nil {
		log.Printf("清理过期验证码失败: %v", result.Error)
	} else {
		log.Printf("清理过期验证码成功，删除了 %d 条记录", result.RowsAffected)
	}
}
