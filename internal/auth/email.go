package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"bagr-backend/internal/utils"
)

// EmailService handles email operations
type EmailService struct {
	clientID     string
	clientSecret string
	tenantID     string
	fromEmail    string
	fromName     string
	testMode     bool // For testing without actual email sending
	accessToken  string
	tokenExpiry  time.Time
}

// EmailConfig represents email configuration
type EmailConfig struct {
	ClientID     string
	ClientSecret string
	TenantID     string
	FromEmail    string
	FromName     string
	TestMode     bool // For testing without actual email sending
}

// NewEmailService creates a new email service
func NewEmailService(config EmailConfig) *EmailService {
	return &EmailService{
		clientID:     config.ClientID,
		clientSecret: config.ClientSecret,
		tenantID:     config.TenantID,
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

// sendEmail sends an email using Microsoft Graph API
func (e *EmailService) sendEmail(to, subject, body string) error {
	logger := utils.GetLogger()

	logger.WithFields(map[string]interface{}{
		"to":         to,
		"subject":    subject,
		"from_email": e.fromEmail,
	}).Debug("Preparing to send email via Microsoft Graph API")

	// Get access token
	token, err := e.getAccessToken()
	if err != nil {
		logger.WithError(err).Error("Failed to get access token")
		return fmt.Errorf("failed to get access token: %w", err)
	}

	// Prepare the email message for Graph API
	emailData := map[string]interface{}{
		"message": map[string]interface{}{
			"subject": subject,
			"body": map[string]interface{}{
				"contentType": "HTML",
				"content":     body,
			},
			"toRecipients": []map[string]interface{}{
				{
					"emailAddress": map[string]string{
						"address": to,
					},
				},
			},
		},
		"saveToSentItems": true,
	}

	jsonData, err := json.Marshal(emailData)
	if err != nil {
		logger.WithError(err).Error("Failed to marshal email data")
		return fmt.Errorf("failed to marshal email data: %w", err)
	}

	// Send email via Graph API
	url := fmt.Sprintf("https://graph.microsoft.com/v1.0/users/%s/sendMail", e.fromEmail)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		logger.WithError(err).Error("Failed to create HTTP request")
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		logger.WithError(err).Error("Failed to send HTTP request to Graph API")
		return fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		body, _ := io.ReadAll(resp.Body)
		logger.WithFields(map[string]interface{}{
			"status_code": resp.StatusCode,
			"response":    string(body),
		}).Error("Graph API returned error")
		return fmt.Errorf("Graph API error: status %d, response: %s", resp.StatusCode, string(body))
	}

	logger.WithFields(map[string]interface{}{
		"to":      to,
		"subject": subject,
	}).Info("Email sent successfully via Microsoft Graph API")
	return nil
}

// getAccessToken gets an access token for Microsoft Graph API
func (e *EmailService) getAccessToken() (string, error) {
	// Check if we have a valid token
	if e.accessToken != "" && time.Now().Before(e.tokenExpiry) {
		return e.accessToken, nil
	}

	logger := utils.GetLogger()
	logger.Debug("Getting new access token from Microsoft Graph API")

	// Prepare token request
	data := url.Values{}
	data.Set("client_id", e.clientID)
	data.Set("client_secret", e.clientSecret)
	data.Set("scope", "https://graph.microsoft.com/.default")
	data.Set("grant_type", "client_credentials")

	req, err := http.NewRequest("POST", fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", e.tenantID), strings.NewReader(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("failed to create token request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send token request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		logger.WithFields(map[string]interface{}{
			"status_code": resp.StatusCode,
			"response":    string(body),
		}).Error("Token request failed")
		return "", fmt.Errorf("token request failed: status %d, response: %s", resp.StatusCode, string(body))
	}

	var tokenResp struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
		TokenType   string `json:"token_type"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", fmt.Errorf("failed to decode token response: %w", err)
	}

	// Store token and expiry
	e.accessToken = tokenResp.AccessToken
	e.tokenExpiry = time.Now().Add(time.Duration(tokenResp.ExpiresIn-60) * time.Second) // 60 seconds buffer

	logger.Debug("Successfully obtained access token from Microsoft Graph API")
	return e.accessToken, nil
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
