package user

type CreateUserInput struct {
	Email    string `validate:"required,email"`
	Name     string `validate:"required,min=3,max=100"`
	Password string `validate:"required,min=8"`
}

type SignInUserInput struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=8"`
}
