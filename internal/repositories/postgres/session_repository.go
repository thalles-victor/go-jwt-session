package postgres

import (
	"context"
	"jwt-session/internal/models"
	"time"

	"github.com/jmoiron/sqlx"
)

// SessionRepository define o contrato para operações de sessão
type SessionRepository interface {
	GetAllByUserID(ctx context.Context, userID string) ([]models.Session, error)
	GetByID(ctx context.Context, id string) (*models.Session, error)
	Create(ctx context.Context, session *models.Session) (*models.Session, error)
	DeleteByID(ctx context.Context, id string) error
	DeleteByUserID(ctx context.Context, userID string) error
}

type sessionRepository struct {
	db *sqlx.DB
}

func NewSessionRepository(db *sqlx.DB) SessionRepository {
	return &sessionRepository{db: db}
}

func (r *sessionRepository) GetAllByUserID(ctx context.Context, userID string) ([]models.Session, error) {
	var sessions []models.Session
	err := r.db.SelectContext(ctx, &sessions, `
		SELECT id, user_id, browser, ip, created_at, expires_at
		FROM sessions
		WHERE user_id = $1
	`, userID)
	if err != nil {
		return nil, err
	}
	return sessions, nil
}

func (r *sessionRepository) GetByID(ctx context.Context, id string) (*models.Session, error) {
	var session models.Session
	err := r.db.GetContext(ctx, &session, `
		SELECT id, user_id, browser, ip, created_at, expires_at
		FROM sessions
		WHERE id = $1
	`, id)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *sessionRepository) Create(ctx context.Context, session *models.Session) (*models.Session, error) {
	if session.CreatedAt.IsZero() {
		session.CreatedAt = time.Now()
	}

	query := `
		INSERT INTO sessions (id, user_id, browser, ip, created_at, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, user_id, browser, ip, created_at, expires_at
	`

	var createdSession models.Session
	err := r.db.GetContext(ctx, &createdSession, query,
		session.ID,
		session.UserID,
		session.Browser,
		session.IP,
		session.CreatedAt,
		session.ExpiresAt,
	)
	if err != nil {
		return nil, err
	}

	return &createdSession, nil
}

func (r *sessionRepository) DeleteByID(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM sessions
		WHERE id = $1
	`, id)
	return err
}

func (r *sessionRepository) DeleteByUserID(ctx context.Context, userID string) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM sessions
		WHERE user_id = $1
	`, userID)
	return err
}
