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

	"gopkg.in/gomail.v2"
)

type EmailService struct{}

func NewEmailService() *EmailService {
	return &EmailService{}
}

type VerificationCode struct {
	Code      string
	Email     string
	ExpiresAt time.Time
}

var verificationCodes = make(map[string]*VerificationCode)

func (s *EmailService) SendVerificationCode(email string) error {
	code := generateVerificationCode()

	cfg := config.AppConfig

	log.Printf("准备发送验证码邮件到: %s", email)
	log.Printf("SMTP 配置: Host=%s, Port=%s, User=%s", cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUser)

	from := cfg.SMTPUser

	to := []string{email}

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

	if port == 465 {
		d.SSL = true
	}

	d.TLSConfig = &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         cfg.SMTPHost,
	}

	log.Printf("开始连接 SMTP 服务器: %s:%d", cfg.SMTPHost, port)

	if err := d.DialAndSend(m); err != nil {
		log.Printf("发送邮件失败: %v", err)
		return fmt.Errorf("发送邮件失败: %v", err)
	}

	log.Printf("邮件发送成功")

	verificationCodes[email] = &VerificationCode{
		Code:      code,
		Email:     email,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}

	return nil
}

func (s *EmailService) SendResetPasswordCode(email string) error {
	code := generateVerificationCode()

	cfg := config.AppConfig

	log.Printf("准备发送密码重置验证码邮件到: %s", email)
	log.Printf("SMTP 配置: Host=%s, Port=%s, User=%s", cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUser)

	from := cfg.SMTPUser

	to := []string{email}

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

	if port == 465 {
		d.SSL = true
	}

	d.TLSConfig = &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         cfg.SMTPHost,
	}

	log.Printf("开始连接 SMTP 服务器: %s:%d", cfg.SMTPHost, port)

	if err := d.DialAndSend(m); err != nil {
		log.Printf("发送邮件失败: %v", err)
		return fmt.Errorf("发送邮件失败: %v", err)
	}

	log.Printf("邮件发送成功")

	verificationCodes[email] = &VerificationCode{
		Code:      code,
		Email:     email,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}

	return nil
}

func (s *EmailService) VerifyCode(email, code string) bool {
	vc, exists := verificationCodes[email]
	if !exists {
		return false
	}

	if vc.Code != code {
		return false
	}

	if time.Now().After(vc.ExpiresAt) {
		delete(verificationCodes, email)
		return false
	}

	delete(verificationCodes, email)
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

func (s *EmailService) CleanupExpiredCodes() {
	for email, vc := range verificationCodes {
		if time.Now().After(vc.ExpiresAt) {
			delete(verificationCodes, email)
		}
	}
}
