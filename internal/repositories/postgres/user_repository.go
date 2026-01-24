package postgres

import (
	"context"
	"jwt-session/internal/models"
	"time"

	"github.com/jmoiron/sqlx"
)

// UserRepository define o contrato para operações de usuário
type UserRepository interface {
	GetByID(ctx context.Context, id string) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Create(ctx context.Context, user *models.User) (*models.User, error)
	UpdatePassword(ctx context.Context, userID string, newPassword string) error
}

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	err := r.db.GetContext(ctx, &user, `SELECT id, name, email, password, created_at, updated_at FROM "users" WHERE id = $1`, id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.db.GetContext(ctx, &user, `SELECT id, name, email, password, created_at, updated_at FROM "users" WHERE email = $1`, email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Create(ctx context.Context, user *models.User) (*models.User, error) {
	if user.CreatedAt.IsZero() {
		user.CreatedAt = time.Now()
	}

	query := `
		INSERT INTO "users" (id, name, email, password, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, name, email, password, created_at, updated_at
	`

	var createdUser models.User
	if err := r.db.GetContext(ctx, &createdUser, query,
		user.ID,
		user.Name,
		user.Email,
		user.Password,
		user.CreatedAt,
	); err != nil {
		return nil, err
	}

	return &createdUser, nil
}

func (r *userRepository) UpdatePassword(ctx context.Context, userID string, newPassword string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE "users"
		SET password = $1, updated_at = $2
		WHERE id = $3
	`, newPassword, time.Now(), userID)

	return err
}
