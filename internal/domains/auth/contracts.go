package auth

import (
	"context"
	"jwt-session/internal/models"
	postgresRepo "jwt-session/internal/repositories/postgres"
	"time"
)

// Service define o contrato principal do domínio de autenticação
type Service interface {
	SignUp(ctx context.Context, dto SignUpDto) (*SignUpResponse, error)
	SignIn(ctx context.Context, dto SignInDto) (*SignInResponse, error)
	SignUpWithSession(ctx context.Context, dto SignUpDto, browser, ip string) (*SignUpWithSessionResponse, error)
	SignInWithSession(ctx context.Context, dto SignInDto, browser, ip string) (*SignInWithSessionResponse, error)
	GetAllSessionsFromUser(ctx context.Context, userID string) ([]models.Session, error)
	DeleteSessionFromUser(ctx context.Context, userID, sessionID string) error
	RequestRecovery(ctx context.Context, email string) error
	ChangePasswordRequestRecovery(ctx context.Context, dto ChangePasswordRequestRecoveryDto) error
}

// JWTService define o contrato para operações JWT
type JWTService interface {
	GenerateJwt(sub string) (string, error)
	ParseJWT(tokenAsString string) (string, error)
}

// MailService define o contrato para envio de emails
type MailService interface {
	SendCreateAccount(name, email string) error
	SendRecoveryRequestEmail(name, email, resetUrl string) error
	SendConfirmationChangePassword(name, email, loginUrl string) error
}

// CodeService define o contrato para geração de códigos
type CodeService interface {
	GenerateRecoveryCode(n int) (string, error)
}

// DateService define o contrato para operações de data
type DateService interface {
	GenerateFutureDate(value int, unit string) (time.Time, error)
	IsNotExpired(expiresAt time.Time) bool
}

// Repositories agrupa os repositórios necessários
type Repositories struct {
	User     postgresRepo.UserRepository
	Session  postgresRepo.SessionRepository
	Recovery postgresRepo.RecoveryRepository
}

// SignUpResponse representa a resposta do SignUp
type SignUpResponse struct {
	User        *models.User `json:"user"`
	AccessToken struct {
		JWT     string `json:"jwt"`
		Expires string `json:"expires"`
	} `json:"access_token"`
}

// SignInResponse representa a resposta do SignIn
type SignInResponse struct {
	User        *models.User `json:"user"`
	AccessToken struct {
		JWT     string `json:"jwt"`
		Expires string `json:"expires"`
	} `json:"access_token"`
}

// SignUpWithSessionResponse representa a resposta do SignUpWithSession
type SignUpWithSessionResponse struct {
	User        *models.User        `json:"user"`
	AccessToken struct {
		JWT     string `json:"jwt"`
		Expires string `json:"expires"`
	} `json:"access_token"`
	Session *models.Session `json:"session"`
}

// SignInWithSessionResponse representa a resposta do SignInWithSession
type SignInWithSessionResponse struct {
	User        *models.User        `json:"user"`
	AccessToken struct {
		JWT     string `json:"jwt"`
		Expires string `json:"expires"`
	} `json:"access_token"`
	Session *models.Session `json:"session"`
}
