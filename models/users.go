package models

import (
	"context"
	"database/sql"
	"strings"

	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type UserLogin struct {
	Email    string
	Password string
}

type ForgetPassword struct {
	Email string
}

type UserRegister struct {
	Email          string
	Password       string
	RepeatPassword string
}

type UserNewPassword struct {
	Password       string
	RepeatPassword string
}

// to use main db that initialised in main.go
type UserModel struct {
	db     *sql.DB
	redis  *redis.Client
	ctx    context.Context
	cancel context.CancelFunc
}

func (m *UserModel) InsertUser(email, password, hashed string) (int64, error) {
	HashedPassword, err := m.GeneratePassword(password)
	if err != nil {
		return 0, err
	}

	result, err := m.db.Exec("INSERT INTO users(`email`,`password`,`activation_token`) VALUES (?, ?,? )", email, string(HashedPassword), hashed)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (m *UserModel) SetLoginToken(token string, uid int64) error {
	_, err := m.db.Exec("UPDATE `users` SET `login_token` = ? WHERE `id` = ?", token, uid)
	return err
}

// logout
func (m *UserModel) Logout(uid int64) error {
	_, err := m.db.Exec("UPDATE `users` SET `login_token` = NULL WHERE `id` = ?", uid)
	if err != nil {
		return err
	}
	return nil
}

func (m *UserModel) Login(creds *UserLogin) (int64, error) {
	var databasePassword string
	var uid int64

	// Begin a transaction (optional)
	tx, err := m.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback() // Rollback if we encounter an error before commit

	// Query for active user by email
	err = tx.QueryRow("SELECT password, id FROM `users` WHERE `active` = 1 AND `email` = ?", strings.TrimSpace(creds.Email)).
		Scan(&databasePassword, &uid)
	if err != nil {
		if err == sql.ErrNoRows {
			// No user found with the given email
			return 0, ErrUserNotFound
		}
		// Other database error
		return 0, err
	}

	// Compare passwords
	err = bcrypt.CompareHashAndPassword([]byte(databasePassword), []byte(creds.Password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			// Incorrect password
			return 0, ErrIncorrectPassword
		}
		// Other bcrypt error
		return 0, err
	}

	// Commit transaction (optional)
	if err := tx.Commit(); err != nil {
		return 0, err
	}

	// Return user ID on successful login
	return uid, nil
}

func (m *UserModel) EmailExist(email string) int64 {
	var uid int64
	_ = m.db.QueryRow("SELECT `id` FROM `users` WHERE  `email` = ?", email).Scan(&uid)
	return uid
}

func (m *UserModel) ValidUser(token string) (int64, error) {
	var id int64
	err := m.db.QueryRow("SELECT `id` FROM `users` WHERE `login_token` = ? ", token).Scan(&id)
	return id, err
}

func (m *UserModel) ValidURI(uri string) bool {
	var exists int64
	query := "SELECT 1 FROM users WHERE activation_token = ? AND active = 0"
	err := m.db.QueryRow(query, uri).Scan(&exists)
	if err != nil {
		return false
	}

	return exists > 0
}

func (m *UserModel) AccountActivate(token string) error {
	_, err := m.db.Exec("UPDATE `users` SET `activation_token` = NULL, `active` = 1 WHERE `activation_token` = ? ", token)
	return err
}

func (m *UserModel) ForgetPassword(uid int64, uri string) error {
	_, _ = m.db.Exec("UPDATE `forget_passw` SET `superseded` = 1 WHERE `uid` = ?", uid)
	_, err := m.db.Exec("INSERT INTO `forget_passw` (`uid`,`uri`,`superseded`) VALUES(?,?,0) ", uid, uri)
	return err
}

func (m *UserModel) ForgetPasswordUri(uri string) (int64, error) {
	var result int64
	err := m.db.QueryRow("SELECT uid FROM `forget_passw` WHERE `uri` = ? AND `superseded` = 0", uri).Scan(&result)
	if err != nil {
		return 0, err
	}

	return result, nil
}

func (m *UserModel) NewPassword(newPassword string, id int64) error {
	newHashedPassword, err := m.GeneratePassword(newPassword)
	if err != nil {
		return err
	}

	stmt := "UPDATE users SET password = ? WHERE id = ?"
	_, err = m.db.Exec(stmt, string(newHashedPassword), id)
	if err != nil {
		return err
	}

	_, _ = m.db.Exec("UPDATE `forget_passw` SET `superseded` =1 WHERE `uid` = ?", id)
	m.ActivityLog("password_changed", id)
	return nil
}

func (m *UserModel) GeneratePassword(newPassword string) ([]byte, error) {
	newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), 12)
	return newHashedPassword, err
}

func (m *UserModel) ActivityLog(activity string, uid int64) {
	_, _ = m.db.Exec("UPDATE `user_log` SET superseded = 1 WHERE activity = ? AND uid = ?", activity, uid)
	_, _ = m.db.Exec("INSERT INTO `user_log` SET  activity = ? , uid = ?, superseded = 0", activity, uid)
}
