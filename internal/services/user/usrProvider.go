package user

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/steephseqq/maximlibs/logger/sl"
	userpb "github.com/steephseqq/maximprotos-user/gen/go/user"
)

func (s *UserService) User(
	ctx context.Context,
	email string,
) (*userpb.UserEntity, error) {
	const op = "user.UserFromUsername"

	log := s.log.With(
		slog.String("op", op),
	)

	log.Info("attempting to get user from email")

	userDB, err := s.usrProvier.User(
		ctx, email,
	)

	if err != nil {
		log.Error("failed to get users from username", sl.Err(err))
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	log.Info("successfully get users from username")
	return &userpb.UserEntity{
		Id: &userDB.ID,
	}, nil
}

func (s *UserService) UsersFromIDs(
	ctx context.Context,
	userIDs []string,
) ([]*userpb.UserEntity, error) {
	const op = "user.UsersFromIDs"

	log := s.log.With(
		slog.String("op", op),
	)

	log.Info("attemping to get users from ids")

	usersDB, err := s.usrProvier.Users(
		ctx, userIDs,
	)

	if err != nil {
		log.Error("failed to get users", sl.Err(err))
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	usersPB := make([]*userpb.UserEntity, 0)
	for i := range usersDB {
		user := &userpb.UserEntity{
			Id:        usersPB[i].Id,
			Name:      usersPB[i].Name,
			AvatarUrl: usersPB[i].AvatarUrl,
		}
		usersPB = append(usersPB, user)
	}

	return usersPB, nil
}

func (s *UserService) UsersFromUsername(
	ctx context.Context,
	username string,
) ([]*userpb.UserEntity, error) {
	const op = "user.UsersFromIDs"

	log := s.log.With(
		slog.String("op", op),
	)

	log.Info("attemping to get users from username")

	usersDB, err := s.usrProvier.UsersFromUsername(
		ctx, username,
	)

	if err != nil {
		log.Error("failed to get users", sl.Err(err))
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	usersPB := make([]*userpb.UserEntity, 0)
	for i := range usersDB {
		user := &userpb.UserEntity{
			Id:        usersPB[i].Id,
			Name:      usersPB[i].Name,
			AvatarUrl: usersPB[i].AvatarUrl,
		}
		usersPB = append(usersPB, user)
	}

	return usersPB, nil
}
