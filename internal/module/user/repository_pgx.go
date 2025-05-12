package user

import (
	"context"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type sqlRepo struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &sqlRepo{db}
}

func (r *sqlRepo) GetCredentialsByEmail(ctx context.Context, email string) (*Credential, error) {
	var credential Credential
	query := `SELECT credentials.user_id, credentials.password_hash, credentials.created_at
              FROM credentials
              JOIN users
              ON credentials.user_id = users.id
              WHERE users.email = $1`
	if err := pgxscan.Get(ctx, r.db, &credential, query, email); err != nil {
		return nil, err
	}
	return &credential, nil
}

func (r *sqlRepo) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	query := `SELECT users.id, email, phone, is_email_verified, is_phone_verified, created_at, updated_at
              FROM users
              WHERE email = $1`
	if err := pgxscan.Get(ctx, r.db, &user, query, email); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *sqlRepo) Create(ctx context.Context, input CreateUser) (*User, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		err := tx.Rollback(ctx)
		if err != nil {

		}
	}(tx, ctx)

	queryCreateUser := `
		INSERT INTO users (email)
		VALUES ($1)
		RETURNING id, email, phone, is_email_verified, is_phone_verified, created_at, updated_at;
	`
	var user User
	if err := pgxscan.Get(ctx, tx, &user, queryCreateUser, input.Email); err != nil {
		return nil, err
	}

	queryPassword := `
		INSERT INTO credentials (user_id, password_hash)
		VALUES ($1, $2)
	`

	if _, err := tx.Exec(ctx, queryPassword, user.ID, input.PasswordHash); err != nil {
		return nil, err
	}

	if tx.Commit(ctx) != nil {
		return nil, err
	}
	return &user, nil
}

func (r *sqlRepo) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`
	if err := r.db.QueryRow(ctx, query, email).Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}
func (r *sqlRepo) ExistsByUUID(ctx context.Context, uuid string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)`
	if err := r.db.QueryRow(ctx, query, uuid).Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}
