package internal

import "context"

// UserRepository provides access a user store.
type UserRepository interface {
	GetUsers(ctx context.Context) ([]*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	CreateUser(ctx context.Context, user *User) error
}

// User is the central class in the domain model.
type User struct {
	ID        int64
	Name      string
	FirstName string
	LastName  string
	Email     string
	Phone     string
	Age       int16
	Position  string
}
