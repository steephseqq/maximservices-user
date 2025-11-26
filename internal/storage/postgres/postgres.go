package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	"user/internal/domain/models"
	"user/internal/storage"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type Storage struct {
	db *sqlx.DB
}

func MustLoad() *Storage {
	connString := os.Getenv("DB_URL")
	if connString == "" {
		panic("connString is required")
	}

	storage, err := New(connString)
	if err != nil {
		panic(err)
	}
	return storage
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.postgres.New"

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("%s:%s", op, "dbURL is required")
	}

	db, err := sqlx.Connect("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveUser(
	ctx context.Context,
	u models.User,
) (err error) {
	const op = "storage.postgres.CreateUser"

	query := `
		INSERT
		INTO users
		(uuid,email,username,name,avatar_url,last_seen)
		VALUES ($1,$2,$3,$4,$5,$6)`

	_, err = s.db.Exec(query,
		u.UUID, u.Email, u.Username, u.Name, u.AvatarURL, time.Now())
	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return fmt.Errorf("%s:%w", op, storage.ErrUserExists)
		}
		return fmt.Errorf("%s:%w", op, err)
	}
	return nil
}

func (s *Storage) RemoveUser(
	ctx context.Context,
	uuid string,
) (err error) {
	const op = "storage.postgres.DeleteUser"

	query := `
		DELETE
		FROM users
		WHERE uuid=$1`

	result, err := s.db.Exec(query, uuid)
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("%s:%w", op, storage.ErrUserNotFound)
	}

	return nil
}

func (s *Storage) User(
	ctx context.Context,
	email string,
) (u models.User, err error) {
	const op = "storage.postgres.User"

	query := `SELECT
		uuid,
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

func (s *Storage) DeleteUser(
	ctx context.Context,
	userID string,
	) (string,error){
	const op = "postgres.DeleteUser"

	query:=`
		DELETE 
		FROM users 
		WHERE id = $1`

	if err:=s.db.Exec(query,userID);err!=nil{
		return "",fmt.Errorf("%s:%w",op,err)	
	}

	return userID,nil
}
