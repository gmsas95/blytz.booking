package email

import (
	"fmt"
	"net/smtp"
	"strings"
)

type EmailService struct {
	from     string
	host     string
	port     int
	username string
	password string
}

type EmailConfig struct {
	From     string
	Host     string
	Port     int
	Username string
	Password string
}

func NewEmailService(config EmailConfig) *EmailService {
	return &EmailService{
		from:     config.From,
		host:     config.Host,
		port:     config.Port,
		username: config.Username,
		password: config.Password,
	}
}

func (e *EmailService) SendPasswordReset(to, name, token string) error {
	subject := "Reset Your Password"
	resetURL := fmt.Sprintf("https://blytz.cloud/reset-password?token=%s", token)

	body := fmt.Sprintf(`
		<h2>Hi %s,</h2>
		<p>You requested a password reset for your Blytz.Cloud account.</p>
		<p>Click the link below to reset your password:</p>
		<p><a href="%s" style="background: #3b82f6; color: white; padding: 12px 24px; text-decoration: none; border-radius: 6px; display: inline-block;">Reset Password</a></p>
		<p>This link will expire in 1 hour.</p>
		<p>If you didn't request this, you can safely ignore this email.</p>
	`, name, resetURL)

	return e.sendEmail(to, subject, body)
}

func (e *EmailService) SendBookingConfirmation(to, name, serviceName, date string, depositPaid float64) error {
	subject := "Booking Confirmed"
	body := fmt.Sprintf(`
		<h2>Hi %s,</h2>
		<p>Your booking has been confirmed!</p>
		<h3>Booking Details</h3>
		<p><strong>Service:</strong> %s</p>
		<p><strong>Date:</strong> %s</p>
		<p><strong>Deposit Paid:</strong> $%.2f</p>
		<p>We're looking forward to seeing you!</p>
	`, name, serviceName, date, depositPaid)

	return e.sendEmail(to, subject, body)
}

func (e *EmailService) SendBookingCancellation(to, name, serviceName, date string) error {
	subject := "Booking Cancelled"
	body := fmt.Sprintf(`
		<h2>Hi %s,</h2>
		<p>Your booking has been cancelled.</p>
		<h3>Cancelled Booking Details</h3>
		<p><strong>Service:</strong> %s</p>
		<p><strong>Originally Scheduled:</strong> %s</p>
		<p>Your deposit has been refunded.</p>
		<p>We hope to see you again soon!</p>
	`, name, serviceName, date)

	return e.sendEmail(to, subject, body)
}

func (e *EmailService) sendEmail(to, subject, body string) error {
	// Build email message
	var message strings.Builder
	message.WriteString(fmt.Sprintf("From: %s\r\n", e.from))
	message.WriteString(fmt.Sprintf("To: %s\r\n", to))
	message.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	message.WriteString("MIME-version: 1.0;\r\n")
	message.WriteString("Content-Type: text/html; charset=\"UTF-8\";\r\n\r\n")
	message.WriteString(body)

	// Connect and send via SMTP
	auth := smtp.PlainAuth("", e.username, e.password, e.host)
	serverAddr := fmt.Sprintf("%s:%d", e.host, e.port)

	return smtp.SendMail(serverAddr, auth, e.from, []string{to}, []byte(message.String()))
}
