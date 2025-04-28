package mailer

import (
	"fmt"
	"net/smtp"

	"auth-service/config"
	"auth-service/internal/logger"
)

type Mailer interface {
	SendConfirmationEmail(to, token string) error
	SendPasswordResetEmail(to, token string) error
}

type smtpMailer struct {
	cfg *config.Config
}

func NewSMTPMailer(cfg *config.Config) Mailer {
	return &smtpMailer{cfg: cfg}
}

func (m *smtpMailer) SendConfirmationEmail(to, token string) error {
	subject := "Подтверждение email"
	confirmationURL := fmt.Sprintf("http://localhost:5173/auth.html?confirm=%s", token)
	body := fmt.Sprintf(`Подтверждение email

Здравствуйте!

Для подтверждения вашего email перейдите по следующей ссылке:
%s

Если вы не регистрировались на нашем сайте, проигнорируйте это письмо.

С уважением,
Команда Kodland Community`, confirmationURL)

	logger.Info("Sending confirmation email to ", to)
	err := m.sendMail(to, subject, body)
	if err != nil {
		logger.Error("Error sending confirmation email to ", to, ": ", err)
	} else {
		logger.Info("Confirmation email sent to ", to)
	}
	return err
}

func (m *smtpMailer) SendPasswordResetEmail(to, token string) error {
	// Формируем ссылку для сброса пароля
	resetLink := fmt.Sprintf("http://localhost:5173/auth.html?token=%s", token)

	// Формируем текстовое сообщение
	message := fmt.Sprintf(`Сброс пароля

Здравствуйте!

Мы получили запрос на сброс пароля для вашего аккаунта.

Для сброса пароля перейдите по следующей ссылке:
%s

Если вы не запрашивали сброс пароля, проигнорируйте это письмо.
Ссылка действительна в течение 24 часов.

С уважением,
Команда Kodland Community`, resetLink)

	logger.Info("Sending password reset email to ", to)
	err := m.sendMail(to, "Сброс пароля", message)
	if err != nil {
		logger.Error("Error sending password reset email to ", to, ": ", err)
	} else {
		logger.Info("Password reset email sent to ", to)
	}
	return err
}

func (m *smtpMailer) sendMail(to, subject, body string) error {
	from := m.cfg.SMTPUsername
	password := m.cfg.SMTPPassword
	host := m.cfg.SMTPHost
	port := m.cfg.SMTPPort

	addr := fmt.Sprintf("%s:%d", host, port)
	logger.Info("Connecting to SMTP server at ", addr)

	auth := smtp.PlainAuth("", from, password, host)

	msg := []byte(
		"To: " + to + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/plain; charset=\"UTF-8\"\r\n\r\n" +
			body + "\r\n",
	)

	err := smtp.SendMail(addr, auth, from, []string{to}, msg)
	if err != nil {
		logger.Error("SMTP SendMail error: ", err)
	}
	return err
}
