package auth

import (
	"fmt"
	"net/http"
)

// DomainError representa um erro do domínio com informações estruturadas
type DomainError struct {
	Code       string // ex: USER_ALREADY_EXISTS
	Message    string // mensagem amigável
	HTTPStatus int    // código HTTP sugerido
	Op         string // operação que falhou (ex: "auth.signup")
	Err        error  // erro original (wrapping)
}

// Error implementa a interface error
func (e *DomainError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap retorna o erro original para permitir unwrapping
func (e *DomainError) Unwrap() error {
	return e.Err
}

// Is permite comparação de erros por código usando errors.Is
func (e *DomainError) Is(target error) bool {
	t, ok := target.(*DomainError)
	if !ok {
		return false
	}
	return e.Code == t.Code
}

// NewDomainError cria um novo DomainError
func NewDomainError(code, message string, httpStatus int, op string, err error) *DomainError {
	return &DomainError{
		Code:       code,
		Message:    message,
		HTTPStatus: httpStatus,
		Op:         op,
		Err:        err,
	}
}

var (
	// ErrUserAlreadyExists é retornado quando o usuário já está cadastrado
	ErrUserAlreadyExists = NewDomainError(
		"USER_ALREADY_EXISTS",
		"usuário já cadastrado",
		http.StatusConflict,
		"auth.signup",
		nil,
	)

	// ErrUserNotFound é retornado quando o usuário não é encontrado
	ErrUserNotFound = NewDomainError(
		"USER_NOT_FOUND",
		"usuário não encontrado",
		http.StatusNotFound,
		"auth.signin",
		nil,
	)

	// ErrInvalidPassword é retornado quando a senha é inválida
	ErrInvalidPassword = NewDomainError(
		"INVALID_PASSWORD",
		"senha inválida",
		http.StatusUnauthorized,
		"auth.signin",
		nil,
	)

	// ErrRecoveryNotFound é retornado quando os dados de recuperação não são encontrados
	ErrRecoveryNotFound = NewDomainError(
		"RECOVERY_NOT_FOUND",
		"dados de recuperação não encontrados",
		http.StatusNotFound,
		"auth.recovery",
		nil,
	)

	// ErrRecoveryExpired é retornado quando o código de recuperação expirou
	ErrRecoveryExpired = NewDomainError(
		"RECOVERY_EXPIRED",
		"código de recuperação expirado",
		http.StatusNotAcceptable,
		"auth.recovery",
		nil,
	)

	// ErrRecoveryAttemptsExceeded é retornado quando o número de tentativas foi excedido
	ErrRecoveryAttemptsExceeded = NewDomainError(
		"RECOVERY_ATTEMPTS_EXCEEDED",
		"número de tentativas excedido",
		http.StatusNotAcceptable,
		"auth.recovery",
		nil,
	)

	// ErrInvalidRecoveryCode é retornado quando o código de recuperação é inválido
	ErrInvalidRecoveryCode = NewDomainError(
		"INVALID_RECOVERY_CODE",
		"código de recuperação inválido",
		http.StatusUnauthorized,
		"auth.recovery",
		nil,
	)

	// ErrSessionNotFound é retornado quando a sessão não é encontrada
	ErrSessionNotFound = NewDomainError(
		"SESSION_NOT_FOUND",
		"sessão não encontrada",
		http.StatusNotFound,
		"auth.session",
		nil,
	)

	// ErrSessionExpired é retornado quando a sessão expirou
	ErrSessionExpired = NewDomainError(
		"SESSION_EXPIRED",
		"sessão expirada",
		http.StatusUnauthorized,
		"auth.session",
		nil,
	)

	// ErrSessionNotOwned é retornado quando a sessão não pertence ao usuário
	ErrSessionNotOwned = NewDomainError(
		"SESSION_NOT_OWNED",
		"sessão não pertence ao usuário",
		http.StatusForbidden,
		"auth.session",
		nil,
	)
)
