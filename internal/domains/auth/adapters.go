package auth

import (
	"jwt-session/internal/shared/code"
	"jwt-session/internal/shared/date"
	"jwt-session/internal/shared/jwt"
	"jwt-session/internal/shared/mail"
	"time"
)

// JWTAdapter implementa a interface JWTService
type JWTAdapter struct{}

func NewJWTAdapter() JWTService {
	return &JWTAdapter{}
}

func (a *JWTAdapter) GenerateJwt(sub string) (string, error) {
	return jwt.GenerateJwt(sub)
}

func (a *JWTAdapter) ParseJWT(tokenAsString string) (string, error) {
	return jwt.ParseJWT(tokenAsString)
}

// MailAdapter implementa a interface MailService
type MailAdapter struct{}

func NewMailAdapter() MailService {
	return &MailAdapter{}
}

func (a *MailAdapter) SendCreateAccount(name, email string) error {
	return mail.SendCreateAccount(name, email)
}

func (a *MailAdapter) SendRecoveryRequestEmail(name, email, resetUrl string) error {
	return mail.SendRecoveryRequestEmail(name, email, resetUrl)
}

func (a *MailAdapter) SendConfirmationChangePassword(name, email, loginUrl string) error {
	return mail.SendConfirmationChangePassword(name, email, loginUrl)
}

// CodeAdapter implementa a interface CodeService
type CodeAdapter struct{}

func NewCodeAdapter() CodeService {
	return &CodeAdapter{}
}

func (a *CodeAdapter) GenerateRecoveryCode(n int) (string, error) {
	return code.GenerateRecoveryCode(n)
}

// DateAdapter implementa a interface DateService
type DateAdapter struct{}

func NewDateAdapter() DateService {
	return &DateAdapter{}
}

func (a *DateAdapter) GenerateFutureDate(value int, unit string) (time.Time, error) {
	return date.GenerateFutureDate(value, unit)
}

func (a *DateAdapter) IsNotExpired(expiresAt time.Time) bool {
	return date.IsNotExpired(expiresAt)
}
