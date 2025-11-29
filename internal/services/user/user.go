package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"
	"user/internal/domain/models"
	"user/internal/storage"

	"github.com/steephseqq/maximlibs/logger/sl"
	userpb "github.com/steephseqq/maximprotos-user/gen/go/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	Users(
		ctx context.Context,
		fields []string,
		userID interface{},
		parameter string,
	) ([]models.User, error)
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

var usersFields = map[string]bool{
	"id":         true,
	"name":       true,
	"username":   true,
	"bio":        true,
	"avatar_url": true,
	"last_seen":  true,
}

var usersParameters = map[string]bool{
	"id":       true,
	"username": true,
}

func (s *UserService) Users(
	ctx context.Context,
	fields []string,
	userID interface{},
	parameter string,
) ([]*userpb.UserEntity, error) {
	const op = "user.Users"

	log := s.log.With(
		slog.String("op", op),
	)

	log.Info("attemping to get users")

	validFields := make([]string, 0, len(fields))
	for _, field := range fields {
		if usersFields[field] {
			validFields = append(validFields, field)
		}
	}

	if !usersParameters[parameter] {
		log.Warn("invalid parameter")
		return nil, fmt.Errorf("%s:%w", op, ErrInvalidParameter)
	}

	if len(validFields) == 0 {
		log.Info("valied fields count = 0")
		return nil, fmt.Errorf("%s:%w", op, ErrInvalidFields)
	}

	usersDB, err := s.usrProvier.Users(ctx, fields, userID, parameter)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s:%w", op, storage.ErrUserNotFound)
		}

		log.Error("failed to get users:", sl.Err(err))
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	usersPB := make([]*userpb.UserEntity, len(usersDB))

	for i := range usersDB {
		user := &userpb.UserEntity{
			Id:        &usersDB[i].ID,
			Username:  &usersDB[i].Username,
			Name:      &usersDB[i].Name,
			Bio:       &usersDB[i].Bio,
			AvatarUrl: &usersDB[i].AvatarURL,
			LastSeen:  &usersDB[i].LastSeen,
		}
		usersPB = append(usersPB, user)
	}

	if len(usersPB) == 0 {
		log.Info("usersPB count = 0")
		return nil, status.Error(codes.InvalidArgument, "invalid argument")
	}

	log.Info("getting users is successfully")
	return usersPB, nil
}
