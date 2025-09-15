package auth

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"time"

	"bagr-backend/internal/utils"
)

// EmailService handles email operations
type EmailService struct {
	smtpHost     string
	smtpPort     string
	smtpUsername string
	smtpPassword string
	fromEmail    string
	fromName     string
	testMode     bool // For testing without actual email sending
}

// EmailConfig represents email configuration
type EmailConfig struct {
	SMTPHost     string
	SMTPPort     string
	SMTPUsername string
	SMTPPassword string
	FromEmail    string
	FromName     string
	TestMode     bool // For testing without actual email sending
}

// NewEmailService creates a new email service
func NewEmailService(config EmailConfig) *EmailService {
	return &EmailService{
		smtpHost:     config.SMTPHost,
		smtpPort:     config.SMTPPort,
		smtpUsername: config.SMTPUsername,
		smtpPassword: config.SMTPPassword,
		fromEmail:    config.FromEmail,
		fromName:     config.FromName,
		testMode:     config.TestMode,
	}
}

// SendVerificationEmail sends email verification email
func (e *EmailService) SendVerificationEmail(to, username, token string) error {
	logger := utils.GetLogger()

	logger.WithFields(map[string]interface{}{
		"to":        to,
		"username":  username,
		"test_mode": e.testMode,
	}).Info("Sending verification email")

	subject := "Verify Your Email - BAGR Auction System"

	data := map[string]interface{}{
		"Username":    username,
		"Token":       token,
		"VerifyURL":   fmt.Sprintf("http://localhost:8080/api/v1/auth/verify?token=%s", token),
		"CurrentYear": time.Now().Year(),
	}

	body, err := e.renderTemplate("verification", data)
	if err != nil {
		logger.WithError(err).Error("Failed to render verification email template")
		return fmt.Errorf("failed to render verification template: %w", err)
	}

	// In test mode, just log the verification details
	if e.testMode {
		logger.WithFields(map[string]interface{}{
			"to":         to,
			"subject":    subject,
			"verify_url": data["VerifyURL"],
			"token":      token,
		}).Info("EMAIL VERIFICATION (TEST MODE) - Email content logged instead of sending")

		fmt.Printf("\n=== EMAIL VERIFICATION (TEST MODE) ===\n")
		fmt.Printf("To: %s\n", to)
		fmt.Printf("Subject: %s\n", subject)
		fmt.Printf("Verification URL: %s\n", data["VerifyURL"])
		fmt.Printf("Token: %s\n", token)
		fmt.Printf("=====================================\n\n")
		return nil
	}

	logger.Debug("Sending actual verification email via SMTP")
	err = e.sendEmail(to, subject, body)
	if err != nil {
		logger.WithError(err).Error("Failed to send verification email via SMTP")
		return err
	}

	logger.Info("Verification email sent successfully via SMTP")
	return nil
}

// SendPasswordResetEmail sends password reset email
func (e *EmailService) SendPasswordResetEmail(to, username, token string) error {
	subject := "Reset Your Password - BAGR Auction System"

	data := map[string]interface{}{
		"Username":    username,
		"Token":       token,
		"ResetURL":    fmt.Sprintf("http://localhost:8080/api/v1/auth/reset-password?token=%s", token),
		"CurrentYear": time.Now().Year(),
	}

	body, err := e.renderTemplate("password_reset", data)
	if err != nil {
		return fmt.Errorf("failed to render password reset template: %w", err)
	}

	// In test mode, just log the reset details
	if e.testMode {
		fmt.Printf("\n=== PASSWORD RESET (TEST MODE) ===\n")
		fmt.Printf("To: %s\n", to)
		fmt.Printf("Subject: %s\n", subject)
		fmt.Printf("Reset URL: %s\n", data["ResetURL"])
		fmt.Printf("Token: %s\n", token)
		fmt.Printf("==================================\n\n")
		return nil
	}

	return e.sendEmail(to, subject, body)
}

// SendWelcomeEmail sends welcome email after successful registration
func (e *EmailService) SendWelcomeEmail(to, username, role string) error {
	subject := "Welcome to BAGR Auction System!"

	data := map[string]interface{}{
		"Username":    username,
		"Role":        role,
		"CurrentYear": time.Now().Year(),
	}

	body, err := e.renderTemplate("welcome", data)
	if err != nil {
		return fmt.Errorf("failed to render welcome template: %w", err)
	}

	// In test mode, just log the welcome details
	if e.testMode {
		fmt.Printf("\n=== WELCOME EMAIL (TEST MODE) ===\n")
		fmt.Printf("To: %s\n", to)
		fmt.Printf("Subject: %s\n", subject)
		fmt.Printf("Username: %s\n", username)
		fmt.Printf("Role: %s\n", role)
		fmt.Printf("===============================\n\n")
		return nil
	}

	return e.sendEmail(to, subject, body)
}

// sendEmail sends an email using SMTP
func (e *EmailService) sendEmail(to, subject, body string) error {
	logger := utils.GetLogger()

	logger.WithFields(map[string]interface{}{
		"to":         to,
		"subject":    subject,
		"smtp_host":  e.smtpHost,
		"smtp_port":  e.smtpPort,
		"from_email": e.fromEmail,
	}).Debug("Preparing to send email via SMTP")

	// Create the email message
	msg := fmt.Sprintf("From: %s <%s>\r\n", e.fromName, e.fromEmail)
	msg += fmt.Sprintf("To: %s\r\n", to)
	msg += fmt.Sprintf("Subject: %s\r\n", subject)
	msg += "MIME-Version: 1.0\r\n"
	msg += "Content-Type: text/html; charset=UTF-8\r\n"
	msg += "\r\n"
	msg += body

	// Set up authentication
	auth := smtp.PlainAuth("", e.smtpUsername, e.smtpPassword, e.smtpHost)

	// Send the email
	addr := fmt.Sprintf("%s:%s", e.smtpHost, e.smtpPort)
	logger.WithField("smtp_addr", addr).Debug("Connecting to SMTP server")

	err := smtp.SendMail(addr, auth, e.fromEmail, []string{to}, []byte(msg))
	if err != nil {
		logger.WithError(err).WithFields(map[string]interface{}{
			"to":        to,
			"subject":   subject,
			"smtp_addr": addr,
		}).Error("Failed to send email via SMTP")
		return fmt.Errorf("failed to send email: %w", err)
	}

	logger.WithFields(map[string]interface{}{
		"to":      to,
		"subject": subject,
	}).Info("Email sent successfully via SMTP")
	return nil
}

// renderTemplate renders an HTML email template
func (e *EmailService) renderTemplate(templateName string, data map[string]interface{}) (string, error) {
	tmpl, err := template.New(templateName).Parse(getEmailTemplate(templateName))
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// getEmailTemplate returns the HTML template for the given template name
func getEmailTemplate(templateName string) string {
	templates := map[string]string{
		"verification": `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Email Verification</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: #6366F1; color: white; padding: 20px; text-align: center; border-radius: 8px 8px 0 0; }
        .content { background: #f9f9f9; padding: 30px; border-radius: 0 0 8px 8px; }
        .button { display: inline-block; background: #6366F1; color: white; padding: 12px 24px; text-decoration: none; border-radius: 6px; margin: 20px 0; }
        .footer { text-align: center; margin-top: 30px; color: #666; font-size: 14px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>BAGR Auction System</h1>
        </div>
        <div class="content">
            <h2>Verify Your Email Address</h2>
            <p>Hello {{.Username}},</p>
            <p>Thank you for registering with BAGR Auction System! To complete your registration, please verify your email address by clicking the button below:</p>
            <p style="text-align: center;">
                <a href="{{.VerifyURL}}" class="button">Verify Email Address</a>
            </p>
            <p>If the button doesn't work, you can copy and paste this link into your browser:</p>
            <p style="word-break: break-all; background: #eee; padding: 10px; border-radius: 4px;">{{.VerifyURL}}</p>
            <p>This link will expire in 24 hours for security reasons.</p>
            <p>If you didn't create an account with us, please ignore this email.</p>
        </div>
        <div class="footer">
            <p>&copy; {{.CurrentYear}} BAGR Auction System. All rights reserved.</p>
        </div>
    </div>
</body>
</html>`,
		"password_reset": `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Password Reset</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: #EF4444; color: white; padding: 20px; text-align: center; border-radius: 8px 8px 0 0; }
        .content { background: #f9f9f9; padding: 30px; border-radius: 0 0 8px 8px; }
        .button { display: inline-block; background: #EF4444; color: white; padding: 12px 24px; text-decoration: none; border-radius: 6px; margin: 20px 0; }
        .footer { text-align: center; margin-top: 30px; color: #666; font-size: 14px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>BAGR Auction System</h1>
        </div>
        <div class="content">
            <h2>Reset Your Password</h2>
            <p>Hello {{.Username}},</p>
            <p>We received a request to reset your password for your BAGR Auction System account. Click the button below to reset your password:</p>
            <p style="text-align: center;">
                <a href="{{.ResetURL}}" class="button">Reset Password</a>
            </p>
            <p>If the button doesn't work, you can copy and paste this link into your browser:</p>
            <p style="word-break: break-all; background: #eee; padding: 10px; border-radius: 4px;">{{.ResetURL}}</p>
            <p>This link will expire in 1 hour for security reasons.</p>
            <p>If you didn't request a password reset, please ignore this email. Your password will remain unchanged.</p>
        </div>
        <div class="footer">
            <p>&copy; {{.CurrentYear}} BAGR Auction System. All rights reserved.</p>
        </div>
    </div>
</body>
</html>`,
		"welcome": `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Welcome to BAGR</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: #10B981; color: white; padding: 20px; text-align: center; border-radius: 8px 8px 0 0; }
        .content { background: #f9f9f9; padding: 30px; border-radius: 0 0 8px 8px; }
        .button { display: inline-block; background: #10B981; color: white; padding: 12px 24px; text-decoration: none; border-radius: 6px; margin: 20px 0; }
        .footer { text-align: center; margin-top: 30px; color: #666; font-size: 14px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Welcome to BAGR!</h1>
        </div>
        <div class="content">
            <h2>Welcome to BAGR Auction System, {{.Username}}!</h2>
            <p>Your account has been successfully created as a <strong>{{.Role}}</strong>.</p>
            <p>You can now:</p>
            <ul>
                <li>Browse live music auctions</li>
                <li>Place bids on your favorite beats</li>
                <li>Connect with producers and artists</li>
                <li>Access exclusive content</li>
            </ul>
            <p style="text-align: center;">
                <a href="http://localhost:8080" class="button">Start Exploring</a>
            </p>
            <p>If you have any questions, feel free to contact our support team.</p>
        </div>
        <div class="footer">
            <p>&copy; {{.CurrentYear}} BAGR Auction System. All rights reserved.</p>
        </div>
    </div>
</body>
</html>`,
	}

	return templates[templateName]
}
