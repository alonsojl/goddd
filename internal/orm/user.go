package orm

import (
	"context"
	"goddd/internal"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserSchema struct {
	ID        int64 `gorm:"primarykey"`
	Name      string
	FirstName string
	LastName  string
	Email     string
	Phone     string
	Age       int16
	UpdatedAt time.Time
	CreatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (UserSchema) TableName() string {
	return "users"
}

type userRepository struct {
	logger *logrus.Logger
	db     *gorm.DB
}

func NewUserRepository(logger *logrus.Logger, db *gorm.DB) internal.UserRepository {
	return &userRepository{
		logger: logger,
		db:     db,
	}
}

func (u *userRepository) GetUsers(ctx context.Context) ([]*internal.User, error) {
	var dbUsers []UserSchema
	err := u.db.WithContext(ctx).Find(&dbUsers).Error
	if err != nil {
		u.logger.Errorf("error get users: %v", err)
		return nil, err
	}

	if len(dbUsers) == 0 {
		u.logger.Info("there are no user records")
		return nil, internal.ErrNoRows
	}

	var users []*internal.User
	for _, dbUser := range dbUsers {
		users = append(users, &internal.User{
			ID:        dbUser.ID,
			Name:      dbUser.Name,
			FirstName: dbUser.FirstName,
			LastName:  dbUser.LastName,
			Email:     dbUser.Email,
			Phone:     dbUser.Phone,
			Age:       dbUser.Age,
		})
	}

	return users, nil
}

func (u *userRepository) GetUserByEmail(ctx context.Context, email string) (*internal.User, error) {
	return &internal.User{}, nil
}

func (u *userRepository) CreateUser(ctx context.Context, user *internal.User) error {
	dbUser := UserSchema{
		Name:      user.Name,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Phone:     user.Phone,
		Age:       user.Age,
	}
	err := u.db.WithContext(ctx).Create(&dbUser).Error
	if err != nil {
		u.logger.Errorf("error create user: %v", err)
		return err
	}
	user.ID = dbUser.ID
	return err
}
