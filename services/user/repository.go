package user

import (
	"backend/helpers"
	"backend/services/auth"
	"backend/types"
	"backend/utils"
	"database/sql"
	"fmt"
	"time"
)

type UserRepoImpl struct {
	db *sql.DB
}

func NewUserRepoImpl(db *sql.DB) *UserRepoImpl {
	return &UserRepoImpl{db: db}
}

func (u *UserRepoImpl) CreateUser(user types.User) error {
	_, err := u.db.Exec(
		"Insert into users(username,email,country,password_hash,xp,rank_id) values(?,?,?,?,?,?)",
		user.Username, user.Email, user.Country, user.PasswordHash, user.XP, user.RankId,
	)
	return err
}

func (u *UserRepoImpl) GetUserByUsername(username string) (*types.User, error) {
	var user types.User
	row := u.db.QueryRow("Select * from users where username = ?", username)
	err := row.Scan(
		&user.Id,
		&user.Username,
		&user.Email,
		&user.Country,
		&user.RankId,
		&user.XP,
		&user.PasswordHash,
		&user.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil

}

func (u *UserRepoImpl) UpdateUserEmail(id int, newEmail string) error {
	if !utils.IsValidEmail(newEmail) {
		return fmt.Errorf("invalid email format")
	}
	// Check if email already exists
	var exists bool
	err := u.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)", newEmail).Scan(&exists)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("email already in use")
	}
	// Generate verification token
	token := utils.GenerateToken()
	// Store in a "pending_email_updates" table
	_, err = u.db.Exec(`
        INSERT INTO pending_email_updates (user_id, new_email, token, created_at)
        VALUES (?, ?, ?, NOW())
        ON DUPLICATE KEY UPDATE new_email = VALUES(new_email), token = VALUES(token), created_at = NOW()
    `, id, newEmail, token)
	if err != nil {
		return err
	}
	return helpers.SendVerificationEmail(newEmail, token)
}

func (u *UserRepoImpl) VerifyEmailUpdate(token string) error {
	var e types.PendingEmailUpdate
	row := u.db.QueryRow("Select * from pending_email_updates where token = ?", token)
	if err := row.Scan(&e.Id, &e.UserId, &e.NewEmail, &e.Token, &e.CreatedAt, &e.ExpiresAt); err != nil {
		return err
	}
	if e.ExpiresAt.Before(time.Now()) {
		return fmt.Errorf("token expired")
	}
	if _, err := u.db.Exec("update users set email = ? where id = ?", e.NewEmail, e.UserId); err != nil {
		return err
	}
	_, err := u.db.Exec("delete from pending_email_updates where id = ?", e.Id)
	return err
}

func (u *UserRepoImpl) UpdateUserCountry(id string, country string) (*types.User, error) {
	_, err := u.db.Exec("Update users set country = ? where id = ?", country, id)
	if err != nil {
		return nil, err
	}
	var user types.User
	row := u.db.QueryRow("Select * from users where id = ?", id)
	if err := row.Scan(
		&user.Id,
		&user.Username,
		&user.Email,
		&user.Country,
		&user.RankId,
		&user.XP,
		&user.PasswordHash,
		&user.CreatedAt,
	); err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *UserRepoImpl) RequestPasswordReset(id int, code string) error {
	now := time.Now()
	_, err := u.db.Exec(`
		Insert into password_reset_tokens(user_id,code,created_at,expires_at,used)
		values (?,?,?,?,?)`, id, code, now, now.Add(time.Minute*15), false,
	)
	return err
}

func (u *UserRepoImpl) ResetPasswordWithCode(id int, code string, newPassword string) error {
	// check if the user_id + code exist && not expired && not used
	var e types.PasswordResetToken
	row := u.db.QueryRow("SELECT * FROM password_reset_tokens WHERE user_id = ? AND code = ?", id, code)
	if err := row.Scan(&e.Id, &e.UserId, &e.Code, &e.CreatedAt, &e.ExpiresAt, &e.Used); err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("invalid code")
		}
		return err
	}
	if e.ExpiresAt.Before(time.Now()) {
		return fmt.Errorf("code expired")
	}
	if e.Used {
		return fmt.Errorf("code already used")
	}
	hashedPassword, err := auth.HashPassword(newPassword)
	if err != nil {
		return err
	}
	_, err = u.db.Exec("UPDATE users SET password_hash = ? WHERE id = ?", hashedPassword, id)
	if err != nil {
		return err
	}
	// Mark the code as used
	_, err = u.db.Exec("UPDATE password_reset_tokens SET used = TRUE WHERE id = ?", e.Id)
	return err
}
