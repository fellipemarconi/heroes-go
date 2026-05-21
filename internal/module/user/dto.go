package user

import (
	"mime/multipart"
	"time"
)

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
	CreatedAt time.Time `json:"created_at"`
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Image     string    `json:"image,omitempty"`
}

type ForgotPasswordInput struct {
	Email string `validate:"required,email"`
}

type ResetPasswordInput struct {
	Token       string `validate:"required"`
	NewPassword string `validate:"required,min=8"`
}

type UpdateProfileImageInput struct {
	File        multipart.File
	FileHeader  *multipart.FileHeader
	ContentType string
}
