package user

import "time"

type CreateUserInput struct {
	Email    string `validate:"required,email"`
	Name     string `validate:"required,min=3,max=100"`
	Password string `validate:"required,min=8"`
}

type SignInUserInput struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=8"`
}

type UpdateUserInput struct {
	Name        string `validate:"omitempty,min=3,max=100"`
	OldPassword string `validate:"required_with=NewPassword,omitempty,min=8"`
	NewPassword string `validate:"required_with=OldPassword,omitempty,min=8"`
}

type User struct {
	CreatedAt time.Time
	ID        string
	Email     string
	Name      string
}
