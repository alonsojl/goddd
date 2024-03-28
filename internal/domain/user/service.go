package user

import (
	"context"
	"goddd/internal"
	"goddd/pkg/errx"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/sirupsen/logrus"
)

var (
	ErrInvalidArgument = errx.NewErrorf(internal.CodeInvalidArgument, "invalid params")
	ErrInvalidEmail    = errx.NewErrorf(internal.CodeInvalidArgument, "invalid email")
)

type Service interface {
	GetUsers(ctx context.Context) ([]*internal.User, error)
	GetUserByEmail(ctx context.Context, email string) (*internal.User, error)
	CreateUser(ctx context.Context, params *Params) (*internal.User, error)
}

type service struct {
	logger         *logrus.Logger
	userRepository internal.UserRepository
}

func NewUserService(logger *logrus.Logger, userRepository internal.UserRepository) Service {
	return &service{
		logger:         logger,
		userRepository: userRepository,
	}
}

func (u *service) GetUsers(ctx context.Context) ([]*internal.User, error) {
	return u.userRepository.GetUsers(ctx)
}

func (u *service) GetUserByEmail(ctx context.Context, email string) (*internal.User, error) {
	err := validation.Validate(email, EmailRules()...)
	if err != nil {
		u.logger.Errorf("invalid user email: %v", err)
		return &internal.User{}, ErrInvalidEmail.SetOrigin(err)
	}

	return u.userRepository.GetUserByEmail(ctx, email)
}

func (u *service) CreateUser(ctx context.Context, params *Params) (*internal.User, error) {
	err := params.Validate()
	if err != nil {
		u.logger.Errorf("invalid user params: %v", err)
		return &internal.User{}, ErrInvalidArgument.SetOrigin(err)
	}

	user := &internal.User{
		Name:      params.Name,
		FirstName: params.FirstName,
		LastName:  params.LastName,
		Age:       params.Age,
		Position:  params.Position,
		Email:     params.Email,
		Phone:     params.Phone,
	}

	err = u.userRepository.CreateUser(ctx, user)
	if err != nil {
		u.logger.Errorf("error creating user: %v", err)
		return &internal.User{}, err
	}

	return user, nil
}
