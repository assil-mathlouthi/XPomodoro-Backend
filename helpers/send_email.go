package helpers

import (
	"backend/config"
	"fmt"
	"net/smtp"
)

func SendEmail(toEmail, subject, htmlBody string) error {
	from := "assilbensaid308@gmail.com"
	password := config.Envs.GmailToken
	to := []string{toEmail}

	headers := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=\"UTF-8\"\r\n\r\n", from, toEmail, subject)
	msg := []byte(headers + htmlBody)

	return smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, password, "smtp.gmail.com"),
		from, to, msg)
}

func SendVerificationEmail(toEmail, token string) error {
	subject := "Verify Your Email"
	body := fmt.Sprintf(`<!DOCTYPE html>
			<html lang="en">
			<head>
			<meta charset="UTF-8" />
			<meta name="viewport" content="width=device-width, initial-scale=1.0" />
			<title>Verify Your Email</title>
			</head>
			<body style="margin:0; padding:0; font-family:Arial, Helvetica, sans-serif; background:#f6f9fc; color:#0f172a;">
			<table role="presentation" width="100%%" cellspacing="0" cellpadding="0" style="background:#f6f9fc; padding:24px 0;">
				<tr>
				<td align="center">
					<table role="presentation" width="600" cellspacing="0" cellpadding="0" style="background:#ffffff; border-radius:12px; box-shadow:0 2px 8px rgba(0,0,0,0.06); overflow:hidden;">
					<tr>
						<td style="background:#0ea5e9; padding:20px 24px; color:#ffffff; font-size:20px; font-weight:700;">XPomodoro</td>
					</tr>
					<tr>
						<td style="padding:28px 24px;">
						<h1 style="margin:0 0 12px; font-size:22px; color:#0f172a;">Verify your email</h1>
						<p style="margin:0 0 20px; font-size:14px; line-height:1.6; color:#334155;">
							Thanks for signing up! Please confirm this email address by clicking the button below.
						</p>
						<div style="text-align:center; margin:28px 0;">
							<a href="http://localhost:8000/api/v1/verify?token=%s" style="display:inline-block; background:#0ea5e9; color:#ffffff; text-decoration:none; padding:12px 20px; border-radius:8px; font-weight:600;">
							Verify Email
							</a>
						</div>
						<p style="margin:20px 0 0; font-size:12px; color:#64748b;">
							If the button doesn't work, copy and paste this link into your browser:<br/>
							<span style="word-break:break-all; color:#0ea5e9;">http://localhost:8080/api/v1/verify?token=%s</span>
						</p>
						</td>
					</tr>
					<tr>
						<td style="background:#f1f5f9; padding:16px 24px; font-size:12px; color:#64748b; text-align:center;">
						© 2025 XPomodoro. All rights reserved.
						</td>
					</tr>
					</table>
				</td>
				</tr>
			</table>
			</body>
			</html>`,
		token,
		token,
	)
	return SendEmail(toEmail, subject, body)
}

func SendPasswordResetCode(toEmail, code string) error {
	subject := "Your XPomodoro Password Reset Code"
	body := fmt.Sprintf(`<!DOCTYPE html>
			<html lang="en">
			<head>
			<meta charset="UTF-8" />
			<meta name="viewport" content="width=device-width, initial-scale=1.0" />
			<title>Password Reset Code</title>
			</head>
			<body style="margin:0; padding:0; font-family:Arial, Helvetica, sans-serif; background:#f6f9fc; color:#0f172a;">
			<table role="presentation" width="100%%" cellspacing="0" cellpadding="0" style="background:#f6f9fc; padding:24px 0;">
				<tr>
				<td align="center">
					<table role="presentation" width="600" cellspacing="0" cellpadding="0" style="background:#ffffff; border-radius:12px; box-shadow:0 2px 8px rgba(0,0,0,0.06); overflow:hidden;">
					<tr>
						<td style="background:#0ea5e9; padding:20px 24px; color:#ffffff; font-size:20px; font-weight:700;">XPomodoro</td>
					</tr>
					<tr>
						<td style="padding:28px 24px;">
						<h1 style="margin:0 0 12px; font-size:22px; color:#0f172a;">Password Reset Code</h1>
						<p style="margin:0 0 20px; font-size:14px; line-height:1.6; color:#334155;">
							Your password reset code is: <strong>%s</strong>
						</p>
						<p style="margin:0 0 20px; font-size:14px; line-height:1.6; color:#334155;">
							If you did not request a password reset, please ignore this email.
						</p>
						</td>
					</tr>
					<tr>
						<td style="background:#f1f5f9; padding:16px 24px; font-size:12px; color:#64748b; text-align:center;">
						© 2025 XPomodoro. All rights reserved.
						</td>
					</tr>
					</table>
				</td>
				</tr>
			</table>
			</body>
			</html>`,
		code,
	)
	return SendEmail(toEmail, subject, body)
}
