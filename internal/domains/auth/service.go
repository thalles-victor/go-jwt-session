package auth

import (
	"context"
	"database/sql"
	"fmt"
	"jwt-session/internal/models"
	"jwt-session/internal/shared/logger"
	"time"

	"github.com/google/uuid"
)

type service struct {
	repos     Repositories
	jwtSvc    JWTService
	mailSvc   MailService
	codeSvc   CodeService
	dateSvc   DateService
}

func NewService(
	repos Repositories,
	jwtSvc JWTService,
	mailSvc MailService,
	codeSvc CodeService,
	dateSvc DateService,
) Service {
	return &service{
		repos:   repos,
		jwtSvc:  jwtSvc,
		mailSvc: mailSvc,
		codeSvc: codeSvc,
		dateSvc: dateSvc,
	}
}

func (s *service) SignUp(ctx context.Context, dto SignUpDto) (*SignUpResponse, error) {
	logger.Info.Printf("check if user with email %s exist", dto.Email)
	
	user, err := s.repos.User.GetByEmail(ctx, dto.Email)
	if err != nil && err != sql.ErrNoRows {
		logger.Error.Printf("error when get user from database. error: %s \n", err.Error())
		return nil, NewDomainError(
			"INTERNAL_ERROR",
			"erro ao buscar usuário no banco de dados",
			500,
			"auth.signup",
			err,
		)
	}

	if err == nil && user != nil {
		logger.Warn.Printf("user with email %s already registered.\n", dto.Email)
		return nil, ErrUserAlreadyExists
	}

	logger.Info.Printf("create a new user\n")

	newUser := &models.User{
		ID:        uuid.New().String(),
		Name:      dto.Name,
		Email:     dto.Email,
		Password:  dto.Password,
		CreatedAt: time.Now(),
	}

	userCreated, err := s.repos.User.Create(ctx, newUser)
	if err != nil {
		logger.Error.Printf("error when create user with email %s \n", newUser.Email)
		return nil, NewDomainError(
			"INTERNAL_ERROR",
			"erro interno ao criar usuário",
			500,
			"auth.signup",
			err,
		)
	}

	jwtToken, err := s.jwtSvc.GenerateJwt(userCreated.ID)
	if err != nil {
		logger.Error.Printf("error when generate jwt. error: %s", err.Error())
		return nil, NewDomainError(
			"INTERNAL_ERROR",
			"erro interno ao gerar token",
			500,
			"auth.signup",
			err,
		)
	}

	go func() {
		if err := s.mailSvc.SendCreateAccount(userCreated.Name, userCreated.Email); err != nil {
			logger.Error.Printf("error when send email to user to confirm the account creation")
		}
	}()

	return &SignUpResponse{
		User: userCreated,
		AccessToken: struct {
			JWT     string `json:"jwt"`
			Expires string `json:"expires"`
		}{
			JWT:     jwtToken,
			Expires: "1h",
		},
	}, nil
}

func (s *service) SignIn(ctx context.Context, dto SignInDto) (*SignInResponse, error) {
	logger.Info.Printf("check if user with email %s exist", dto.Email)
	
	user, err := s.repos.User.GetByEmail(ctx, dto.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Warn.Printf("user with email %s unregistered. \n", dto.Email)
			return nil, ErrUserNotFound
		}
		logger.Error.Printf("error when get user from database. error: %s \n", err.Error())
		return nil, NewDomainError(
			"INTERNAL_ERROR",
			"erro ao buscar usuário no banco de dados",
			500,
			"auth.signin",
			err,
		)
	}

	logger.Info.Println("check if password is valid")
	if user.Password != dto.Password {
		logger.Warn.Printf("password of user %s is invalid", user.Email)
		return nil, ErrInvalidPassword
	}

	jwtToken, err := s.jwtSvc.GenerateJwt(user.ID)
	if err != nil {
		logger.Error.Printf("error when generate jwt. error: %s", err.Error())
		return nil, NewDomainError(
			"INTERNAL_ERROR",
			"erro interno ao gerar token",
			500,
			"auth.signin",
			err,
		)
	}

	return &SignInResponse{
		User: user,
		AccessToken: struct {
			JWT     string `json:"jwt"`
			Expires string `json:"expires"`
		}{
			JWT:     jwtToken,
			Expires: "1h",
		},
	}, nil
}

func (s *service) SignUpWithSession(ctx context.Context, dto SignUpDto, browser, ip string) (*SignUpWithSessionResponse, error) {
	logger.Info.Printf("check if user with email %s exist", dto.Email)
	
	user, err := s.repos.User.GetByEmail(ctx, dto.Email)
	if err != nil && err != sql.ErrNoRows {
		logger.Error.Printf("error when get user from database. error: %s \n", err.Error())
		return nil, NewDomainError(
			"INTERNAL_ERROR",
			"erro ao buscar usuário no banco de dados",
			500,
			"auth.signupWithSession",
			err,
		)
	}

	if err == nil && user != nil {
		logger.Warn.Printf("user with email %s already registered.\n", dto.Email)
		return nil, ErrUserAlreadyExists
	}

	newUser := &models.User{
		ID:        uuid.New().String(),
		Name:      dto.Name,
		Email:     dto.Email,
		Password:  dto.Password,
		CreatedAt: time.Now(),
	}

	userCreated, err := s.repos.User.Create(ctx, newUser)
	if err != nil {
		logger.Error.Printf("error when create user with email %s \n", newUser.Email)
		return nil, NewDomainError(
			"INTERNAL_ERROR",
			"erro interno ao criar usuário",
			500,
			"auth.signupWithSession",
			err,
		)
	}

	expiresAt, err := s.dateSvc.GenerateFutureDate(1, "days")
	if err != nil {
		logger.Error.Printf("error when generate expiration data to session. error: %s \n", err.Error())
		return nil, NewDomainError(
			"INTERNAL_ERROR",
			"erro interno ao gerar data de expiração da sessão",
			500,
			"auth.signupWithSession",
			err,
		)
	}

	session := &models.Session{
		ID:        uuid.New().String(),
		UserID:    userCreated.ID,
		Browser:   &browser,
		IP:        &ip,
		CreatedAt: time.Now(),
		ExpiresAt: expiresAt,
	}

	sessionCreated, err := s.repos.Session.Create(ctx, session)
	if err != nil {
		logger.Error.Printf("internal server error when create session: error: %s", err.Error())
		return nil, NewDomainError(
			"INTERNAL_ERROR",
			"erro interno ao criar sessão",
			500,
			"auth.signupWithSession",
			err,
		)
	}

	jwtToken, err := s.jwtSvc.GenerateJwt(sessionCreated.ID)
	if err != nil {
		logger.Error.Printf("internal server error when generate jwt. error: %s", err.Error())
		return nil, NewDomainError(
			"INTERNAL_ERROR",
			"erro interno ao gerar token",
			500,
			"auth.signupWithSession",
			err,
		)
	}

	go func() {
		if err := s.mailSvc.SendCreateAccount(userCreated.Name, userCreated.Email); err != nil {
			logger.Error.Printf("error when send email to user to confirm the account creation")
		}
	}()

	return &SignUpWithSessionResponse{
		User: userCreated,
		AccessToken: struct {
			JWT     string `json:"jwt"`
			Expires string `json:"expires"`
		}{
			JWT:     jwtToken,
			Expires: "1h",
		},
		Session: sessionCreated,
	}, nil
}

func (s *service) SignInWithSession(ctx context.Context, dto SignInDto, browser, ip string) (*SignInWithSessionResponse, error) {
	logger.Info.Printf("check if user with email %s exist", dto.Email)
	
	user, err := s.repos.User.GetByEmail(ctx, dto.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Warn.Printf("user with email %s unregistered. \n", dto.Email)
			return nil, ErrUserNotFound
		}
		logger.Error.Printf("error when get user from database. error: %s \n", err.Error())
		return nil, NewDomainError(
			"INTERNAL_ERROR",
			"erro ao buscar usuário no banco de dados",
			500,
			"auth.signinWithSession",
			err,
		)
	}

	logger.Info.Println("check if password is valid")
	if user.Password != dto.Password {
		logger.Warn.Printf("password of user %s is invalid", user.Email)
		return nil, ErrInvalidPassword
	}

	expiresAt, err := s.dateSvc.GenerateFutureDate(1, "days")
	if err != nil {
		logger.Error.Printf("error when generate expiration data to session. error: %s \n", err.Error())
		return nil, NewDomainError(
			"INTERNAL_ERROR",
			"erro interno ao gerar data de expiração da sessão",
			500,
			"auth.signinWithSession",
			err,
		)
	}

	session := &models.Session{
		ID:        uuid.New().String(),
		UserID:    user.ID,
		Browser:   &browser,
		IP:        &ip,
		CreatedAt: time.Now(),
		ExpiresAt: expiresAt,
	}

	sessionCreated, err := s.repos.Session.Create(ctx, session)
	if err != nil {
		logger.Error.Printf("internal server error when create session: error: %s", err.Error())
		return nil, NewDomainError(
			"INTERNAL_ERROR",
			"erro interno ao criar sessão",
			500,
			"auth.signinWithSession",
			err,
		)
	}

	jwtToken, err := s.jwtSvc.GenerateJwt(sessionCreated.ID)
	if err != nil {
		logger.Error.Printf("internal server error when generate jwt. error: %s", err.Error())
		return nil, NewDomainError(
			"INTERNAL_ERROR",
			"erro interno ao gerar token",
			500,
			"auth.signinWithSession",
			err,
		)
	}

	return &SignInWithSessionResponse{
		User: user,
		AccessToken: struct {
			JWT     string `json:"jwt"`
			Expires string `json:"expires"`
		}{
			JWT:     jwtToken,
			Expires: "1h",
		},
		Session: sessionCreated,
	}, nil
}

func (s *service) GetAllSessionsFromUser(ctx context.Context, userID string) ([]models.Session, error) {
	_, err := s.repos.User.GetByID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, NewDomainError(
			"INTERNAL_ERROR",
			"erro interno ao buscar usuário",
			500,
			"auth.getAllSessionsFromUser",
			err,
		)
	}

	sessions, err := s.repos.Session.GetAllByUserID(ctx, userID)
	if err != nil {
		logger.Error.Printf("internal server error when get all sessions from user. error: %s ", err.Error())
		return nil, NewDomainError(
			"INTERNAL_ERROR",
			"erro interno ao buscar sessões",
			500,
			"auth.getAllSessionsFromUser",
			err,
		)
	}

	return sessions, nil
}

func (s *service) DeleteSessionFromUser(ctx context.Context, userID, sessionID string) error {
	_, err := s.repos.User.GetByID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrUserNotFound
		}
		return NewDomainError(
			"INTERNAL_ERROR",
			"erro interno ao buscar usuário",
			500,
			"auth.deleteSessionFromUser",
			err,
		)
	}

	session, err := s.repos.Session.GetByID(ctx, sessionID)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrSessionNotFound
		}
		return NewDomainError(
			"INTERNAL_ERROR",
			"erro interno ao buscar sessão",
			500,
			"auth.deleteSessionFromUser",
			err,
		)
	}

	if session.UserID != userID {
		logger.Warn.Printf("the session with id %s belong another user", sessionID)
		return ErrSessionNotOwned
	}

	if err := s.repos.Session.DeleteByID(ctx, sessionID); err != nil {
		logger.Error.Printf("internal server error when delete session of user. error: %s ", err.Error())
		return NewDomainError(
			"INTERNAL_ERROR",
			"erro interno ao deletar sessão",
			500,
			"auth.deleteSessionFromUser",
			err,
		)
	}

	return nil
}

func (s *service) RequestRecovery(ctx context.Context, email string) error {
	logger.Info.Printf("check if user with email %s exist", email)
	
	user, err := s.repos.User.GetByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Warn.Printf("user with email %s not found", email)
			return ErrUserNotFound
		}
		logger.Error.Printf("error when get user from database. error: %s \n", err.Error())
		return NewDomainError(
			"INTERNAL_ERROR",
			"erro interno ao buscar usuário",
			500,
			"auth.requestRecovery",
			err,
		)
	}

	recoveryCode, err := s.codeSvc.GenerateRecoveryCode(120)
	if err != nil {
		return NewDomainError(
			"INTERNAL_ERROR",
			"erro interno ao gerar código de recuperação",
			500,
			"auth.requestRecovery",
			err,
		)
	}

	expiresAt, err := s.dateSvc.GenerateFutureDate(5, "minutes")
	if err != nil {
		logger.Error.Printf("erro ao gerar data futura %s \n", err.Error())
		return NewDomainError(
			"INTERNAL_ERROR",
			"erro interno ao gerar data futura para o código",
			500,
			"auth.requestRecovery",
			err,
		)
	}

	logger.Info.Printf("check if recovery exist by email %s \n", email)
	recovery, err := s.repos.Recovery.GetByEmail(ctx, email)
	if err != nil && err != sql.ErrNoRows {
		logger.Error.Printf("internal server error when get recovery data")
		return NewDomainError(
			"INTERNAL_ERROR",
			"erro interno ao buscar os dados de recuperação",
			500,
			"auth.requestRecovery",
			err,
		)
	}

	logger.Info.Println("save recovery data in the database")
	if err == sql.ErrNoRows {
		logger.Info.Println("recovery data not found, then create")

		_, err := s.repos.Recovery.Create(ctx, &models.Recovery{
			ID:        0,
			UserID:    user.ID,
			Email:     user.Email,
			Code:      recoveryCode,
			Attempts:  0,
			ExpiresAt: expiresAt,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Expired:   false,
		})
		if err != nil {
			logger.Error.Printf("error when save recovery data in the database. error: %s \n", err.Error())
			return NewDomainError(
				"INTERNAL_ERROR",
				"erro interno ao salvar os dados de recuperação",
				500,
				"auth.requestRecovery",
				err,
			)
		}
	} else {
		logger.Info.Println("data of recovery already exist, then update to new data.")
		if err = s.repos.Recovery.UpdateRecovery(ctx, recovery.ID, recoveryCode, 0, expiresAt, false); err != nil {
			logger.Error.Printf("internal server error when update recovery table %s \n", err.Error())
			return NewDomainError(
				"INTERNAL_ERROR",
				"erro ao atualizar os dados de recuperação",
				500,
				"auth.requestRecovery",
				err,
			)
		}
	}

	recoveryUrl := fmt.Sprintf("https://dominio-do-front/recuperar/tocar-senha?code=%s", recoveryCode)

	go func() {
		if err := s.mailSvc.SendRecoveryRequestEmail(user.Name, user.Email, recoveryUrl); err != nil {
			logger.Error.Printf("error when send recovery email")
		}
	}()

	return nil
}

func (s *service) ChangePasswordRequestRecovery(ctx context.Context, dto ChangePasswordRequestRecoveryDto) error {
	logger.Info.Printf("check if user with email %s exist", dto.Email)
	
	user, err := s.repos.User.GetByEmail(ctx, dto.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrUserNotFound
		}
		logger.Error.Printf("error when get user from database. error: %s \n", err.Error())
		return NewDomainError(
			"INTERNAL_ERROR",
			"erro interno ao buscar usuário",
			500,
			"auth.changePasswordRequestRecovery",
			err,
		)
	}

	recovery, err := s.repos.Recovery.GetByEmail(ctx, dto.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrRecoveryNotFound
		}
		logger.Error.Printf("internal server error when get recovery data: %s \n", err.Error())
		return NewDomainError(
			"INTERNAL_ERROR",
			"erro interno ao buscar dados de recuperação",
			500,
			"auth.changePasswordRequestRecovery",
			err,
		)
	}

	if recovery.Expired {
		logger.Warn.Printf("recovery code already expires")
		return ErrRecoveryExpired
	}

	if recovery.Attempts > 10 {
		logger.Warn.Printf("number of attempts exceed, clear recovery data")
		if err = s.repos.Recovery.MarkAsExpired(ctx, recovery.ID); err != nil {
			logger.Error.Printf("internal server error when make make as expired. error %s", err.Error())
			return NewDomainError(
				"INTERNAL_ERROR",
				"erro interno ao expirar o token",
				500,
				"auth.changePasswordRequestRecovery",
				err,
			)
		}
		return ErrRecoveryAttemptsExceeded
	}

	if !s.dateSvc.IsNotExpired(recovery.ExpiresAt) {
		if err = s.repos.Recovery.MarkAsExpired(ctx, recovery.ID); err != nil {
			logger.Error.Printf("internal server error when disable recovery code")
			return NewDomainError(
				"INTERNAL_ERROR",
				"erro interno ao desativar o código de recuperação",
				500,
				"auth.changePasswordRequestRecovery",
				err,
			)
		}

		logger.Warn.Printf("recovery code already expires at: %s", recovery.ExpiresAt)
		return ErrRecoveryExpired
	}

	logger.Info.Println("check if the code are same")
	if dto.Code != recovery.Code {
		logger.Warn.Println("the recovery are differents. Increasing the attempts.")
		if err := s.repos.Recovery.IncrementAttempts(ctx, recovery.ID); err != nil {
			logger.Error.Printf("error when increments attempts in recovery table. error: %s \n", err.Error())
			return NewDomainError(
				"INTERNAL_ERROR",
				"erro interno ao incrementar tentativas",
				500,
				"auth.changePasswordRequestRecovery",
				err,
			)
		}

		return NewDomainError(
			ErrInvalidRecoveryCode.Code,
			fmt.Sprintf("o código informado é inválido, restam %d tentativas", 10-recovery.Attempts+1),
			ErrInvalidRecoveryCode.HTTPStatus,
			ErrInvalidRecoveryCode.Op,
			nil,
		)
	}

	logger.Info.Printf("code valid to user %s. Update password", user.Email)
	if err := s.repos.User.UpdatePassword(ctx, user.ID, dto.NewPassword); err != nil {
		logger.Error.Printf("internal server error when update user password. error: %s \n", err.Error())
		return NewDomainError(
			"INTERNAL_ERROR",
			"erro interno ao atualizar a senha do usuário",
			500,
			"auth.changePasswordRequestRecovery",
			err,
		)
	}

	logger.Info.Printf("deleting all sessions from use with id: %s", user.ID)
	if err = s.repos.Session.DeleteByUserID(ctx, user.ID); err != nil {
		logger.Error.Printf("internal server error when delete all sessions from user. error: %s", err.Error())
		return NewDomainError(
			"INTERNAL_ERROR",
			"erro interno ao deletar todas as sessões de um usuário",
			500,
			"auth.changePasswordRequestRecovery",
			err,
		)
	}

	logger.Info.Println("disable recovery code")
	if err = s.repos.Recovery.MarkAsExpired(ctx, recovery.ID); err != nil {
		logger.Error.Printf("internal server error when disable recovery code")
		return NewDomainError(
			"INTERNAL_ERROR",
			"erro interno ao desativar o código de recuperação",
			500,
			"auth.changePasswordRequestRecovery",
			err,
		)
	}

	go func() {
		if err := s.mailSvc.SendConfirmationChangePassword(user.Name, user.Email, "http://dominio-do-front-que-tem-que-terlogi"); err != nil {
			logger.Error.Printf("internal server error when send change password confirmation. error: %s", err.Error())
		}
	}()

	return nil
}
