package models

type User struct {
	ID           string `json:"id" db:"id"`
	Email        string `json:"email" db:"email"`
	FirstName    string `json:"firstname" db:"first_name"`
	LastName     string `json:"lastname" db:"last_name"`
	PasswordHash string `json:"password_hash,omitempty" db:"password_hash"`
	Created_At   string `json:"created_at,omitempty" db:"created_at"`
}

type AuthPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=20"`
}

type RegisterPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=20"`
	FistName string `json:"firstname" validate:"required"`
	LastName string `json:"lastname" validate:"required"`
}

type Token struct {
	AccessToken string `json:"access_token"`
}

type UserResponse struct {
	Token Token `json:"token"`
	User  User  `json:"user"`
}

type SessionResponse struct {
	User User `json:"user"`
}
