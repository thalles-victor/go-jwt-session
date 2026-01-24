package postgres

import (
	"context"
	"jwt-session/internal/models"
	"time"

	"github.com/jmoiron/sqlx"
)

// RecoveryRepository define o contrato para operações de recuperação
type RecoveryRepository interface {
	GetByUserID(ctx context.Context, userID string) (*models.Recovery, error)
	GetByEmail(ctx context.Context, email string) (*models.Recovery, error)
	GetByCode(ctx context.Context, code string) (*models.Recovery, error)
	Create(ctx context.Context, recovery *models.Recovery) (*models.Recovery, error)
	IncrementAttempts(ctx context.Context, id int) error
	UpdateRecovery(ctx context.Context, id int, code string, attempts int, expiresAt time.Time, expired bool) error
	MarkAsExpired(ctx context.Context, id int) error
}

type recoveryRepository struct {
	db *sqlx.DB
}

func NewRecoveryRepository(db *sqlx.DB) RecoveryRepository {
	return &recoveryRepository{db: db}
}

func (r *recoveryRepository) GetByUserID(ctx context.Context, userID string) (*models.Recovery, error) {
	var recovery models.Recovery
	err := r.db.GetContext(ctx, &recovery, `
		SELECT id, user_id, email, code, attempts, expires_at, created_at, updated_at, expired
		FROM recoveries
		WHERE user_id = $1
	`, userID)
	if err != nil {
		return nil, err
	}
	return &recovery, nil
}

func (r *recoveryRepository) GetByEmail(ctx context.Context, email string) (*models.Recovery, error) {
	var recovery models.Recovery
	err := r.db.GetContext(ctx, &recovery, `
		SELECT id, user_id, email, code, attempts, expires_at, created_at, updated_at, expired
		FROM recoveries
		WHERE email = $1
	`, email)
	if err != nil {
		return nil, err
	}
	return &recovery, nil
}

func (r *recoveryRepository) GetByCode(ctx context.Context, code string) (*models.Recovery, error) {
	var recovery models.Recovery
	err := r.db.GetContext(ctx, &recovery, `
		SELECT id, user_id, email, code, attempts, expires_at, created_at, updated_at, expired
		FROM recoveries
		WHERE code = $1
	`, code)
	if err != nil {
		return nil, err
	}
	return &recovery, nil
}

func (r *recoveryRepository) Create(ctx context.Context, recovery *models.Recovery) (*models.Recovery, error) {
	if recovery.CreatedAt.IsZero() {
		recovery.CreatedAt = time.Now()
	}
	recovery.UpdatedAt = time.Now()

	query := `
		INSERT INTO recoveries (user_id, email, code, attempts, expires_at, created_at, updated_at, expired)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, user_id, email, code, attempts, expires_at, created_at, updated_at, expired
	`

	var createdRecovery models.Recovery
	if err := r.db.GetContext(ctx, &createdRecovery, query,
		recovery.UserID,
		recovery.Email,
		recovery.Code,
		recovery.Attempts,
		recovery.ExpiresAt,
		recovery.CreatedAt,
		recovery.UpdatedAt,
		recovery.Expired,
	); err != nil {
		return nil, err
	}

	return &createdRecovery, nil
}

func (r *recoveryRepository) IncrementAttempts(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE recoveries
		SET attempts = attempts + 1, updated_at = $2
		WHERE id = $1
	`, id, time.Now())
	return err
}

func (r *recoveryRepository) UpdateRecovery(ctx context.Context, id int, code string, attempts int, expiresAt time.Time, expired bool) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE recoveries
		SET code = $1,
			attempts = $2,
			expires_at = $3,
			expired = $4,
			updated_at = $5
		WHERE id = $6
	`, code, attempts, expiresAt, expired, time.Now(), id)
	return err
}

func (r *recoveryRepository) MarkAsExpired(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE recoveries
		SET expired = TRUE,
			updated_at = $2
		WHERE id = $1
	`, id, time.Now())
	return err
}
