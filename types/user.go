package types

import (
	"time"
)

type UserRepo interface {
	GetUserByUsername(string) (*User, error)
	CreateUser(User) error
	UpdateUserEmail(id int, newEmail string) error
	VerifyEmailUpdate(token string) error
	UpdateUserCountry(id string, country string) (*User, error)
	RequestPasswordReset(id int, code string) error
	ResetPasswordWithCode(id int, code string, newPassword string) error
}

type User struct {
	Id           int       `json:"id"`
	Username     string    `json:"username"`
	Email        *string   `json:"email"`
	PasswordHash string    `json:"password_hash"`
	Country      *string   `json:"country"`
	XP           int       `json:"xp"`
	RankId       int       `json:"rank_id"`
	CreatedAt    time.Time `json:"created_at"`
}

type AuthPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UpdateEmailPayload struct {
	UserId   string `json:"user_id"`
	NewEmail string `json:"new_email"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

type PendingEmailUpdate struct {
	Id        int       `json:"id"`
	UserId    int       `json:"user_id"`
	NewEmail  string    `json:"new_email"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

type UserInfoUpdate struct {
	Country string `json:"country"`
}

type ForgotPasswordPayload struct {
	Username string `json:"username"`
}
type PasswordResetToken struct {
	Id        int       `json:"id"`
	UserId    int       `json:"user_id"`
	Code      string    `json:"code"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
	Used      bool      `json:"used"`
}
type ResetPasswordPayload struct {
	Username    string `json:"username"`
	Code        string `json:"code"`
	NewPassword string `json:"new_password"`
}
