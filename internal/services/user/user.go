package user

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"
	"user/internal/domain/models"

	"github.com/steephseqq/maximlibs/logger/sl"
)

type UserService struct {
	log        *slog.Logger
	usrCreator UserCreator
	usrDeleter UserDeleter
	usrProvier UserProvider
	tokenTTL   time.Duration
}

var (
	ErrInvalidFields      = errors.New("invalid fields")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidParameter   = errors.New("invalid parameter")
)

type UserCreator interface {
	SaveUser(
		ctx context.Context,
		u models.User,
	) (err error)
}

type UserDeleter interface {
	RemoveUser(
		ctx context.Context,
		uuid string,
	) (err error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
	Users(ctx context.Context, userIDs []string) ([]models.User, error)
	UsersFromUsername(ctx context.Context, username string) ([]models.User, error)
}

func New(
	log *slog.Logger,
	usrCreator UserCreator,
	usrDeleter UserDeleter,
	usrProvider UserProvider,
	tokenTTL time.Duration,
) *UserService {
	return &UserService{
		log:        log,
		usrCreator: usrCreator,
		usrProvier: usrProvider,
		tokenTTL:   tokenTTL,
	}
}

func (s *UserService) CreateUser(
	ctx context.Context,
	username, name, email, uuid, avatarURL string,
) (err error) {
	const op = "services.user.CreateUser"

	log := s.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	user := models.User{
		ID:        uuid,
		Email:     email,
		Username:  username,
		Name:      name,
		AvatarURL: avatarURL,
	}
	if err = s.usrCreator.SaveUser(ctx, user); err != nil {
		log.Error("failed to save user", sl.Err(err))
		return fmt.Errorf("%s:%w", op, err)
	}
	return nil
}

func (s *UserService) DeleteUser(
	ctx context.Context,
	uuid string,
) (err error) {
	const op = "services.user.DeleteUser"

	log := s.log.With(
		slog.String("op", op),
	)

	log.Info("attemping to delete user")

	if err := s.usrDeleter.RemoveUser(ctx, uuid); err != nil {
		log.Error("failed to delete user", sl.Err(err))
		return err
	}
	return nil
}
