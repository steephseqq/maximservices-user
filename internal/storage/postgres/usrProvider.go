package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"user/internal/domain/models"
	"user/internal/storage"
)

func (s *Storage) User(
	ctx context.Context,
	email string,
) (u models.User, err error) {
	const op = "storage.postgres.User"

	query := `SELECT
		id,
		username,
		pass_hash
		FROM users
		WHERE email=$1
		`

	var user models.User
	if err := s.db.Select(&user, query, email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s:%w", op, storage.ErrUserNotFound)
		}
		return models.User{}, fmt.Errorf("%s:%w", op, err)
	}
	return user, nil
}

func (s *Storage) Users(
	ctx context.Context,
	userIDs []string,
) ([]models.User, error) {
	const op = "postgres.Users"

	query := `
		SELECT 
			id,
			avatar_url,
			name
		FROM users
		WHERE id = ANY($1)`

	var users []models.User

	if err := s.db.Select(&users, query, userIDs); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return users, nil
		}
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	return users, nil
}

func (s *Storage) UsersFromUsername(
	ctx context.Context,
	username string,
) ([]models.User, error) {
	const op = "postgres.UsersFromUsername"

	username = "%" + username + "%"

	query := `
		SELECT 
			id,
			name,
			avatar_url
		FROM users
		WHERE
			username LIKE $1`

	var users []models.User

	if err := s.db.Select(&users, query, username); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return users, nil
		}
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	return users, nil
}
