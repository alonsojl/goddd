package mysql

import (
	"context"
	"database/sql"
	"errors"
	"goddd/internal"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type UserSchema struct {
	ID        int64  `db:"id"`
	Name      string `db:"name"`
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	Email     string `db:"email"`
	Phone     string `db:"phone"`
	Age       int16  `db:"age"`
}

type userRepository struct {
	logger *logrus.Logger
	db     *sqlx.DB
}

func NewUserRepository(logger *logrus.Logger, db *sqlx.DB) internal.UserRepository {
	return &userRepository{
		logger: logger,
		db:     db,
	}
}

func (u *userRepository) GetUsers(ctx context.Context) ([]*internal.User, error) {
	var (
		query  = `SELECT id, name, first_name, last_name, email, phone, age FROM user`
		dbUser UserSchema
		users  []*internal.User
	)

	rows, err := u.db.Queryx(query)
	if err != nil {
		u.logger.Errorf("error getting users: %v", err)
		return nil, err
	}

	for rows.Next() {
		if err := rows.StructScan(&dbUser); err != nil {
			u.logger.Errorf("error scanning user: %v", err)
			return nil, err
		}
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

	if len(users) == 0 {
		u.logger.Info("there are no user records")
		return nil, internal.ErrNoRows
	}

	return users, nil
}

func (u *userRepository) GetUserByEmail(ctx context.Context, email string) (*internal.User, error) {
	var (
		query  = `SELECT id, name, first_name, last_name, email, phone, age FROM user WHERE email=?`
		dbUser UserSchema
		user   internal.User
	)

	err := u.db.QueryRowxContext(ctx, query, email).StructScan(&dbUser)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			u.logger.Errorf("error getting user: %v", err)
			return nil, err
		}

		return &user, internal.ErrNoRows
	}

	user = internal.User{
		ID:        dbUser.ID,
		Name:      dbUser.Name,
		FirstName: dbUser.FirstName,
		LastName:  dbUser.LastName,
		Email:     dbUser.Email,
		Phone:     dbUser.Phone,
		Age:       dbUser.Age,
	}

	return &user, nil
}

func (u *userRepository) CreateUser(ctx context.Context, user *internal.User) error {
	query := `INSERT INTO user (name, first_name, last_name, email, phone, age) VALUES (?,?,?,?,?,?)`
	result, err := u.db.ExecContext(ctx, query,
		user.Name,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Phone,
		user.Age,
	)
	if err != nil {
		u.logger.Errorf("error creating user: %v", err)
		return err
	}

	user.ID, err = result.LastInsertId()
	if err != nil {
		u.logger.Errorf("error getting user insertion ID: %v", err)
		return err
	}

	return nil
}
