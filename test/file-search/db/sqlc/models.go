// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"time"
)

type User struct {
	ID           int32     `json:"id"`
	Email        string    `json:"email"`
	Username     string    `json:"username"`
	Password     string    `json:"password"`
	PasswordHash string    `json:"password_hash"`
	Phone        string    `json:"phone"`
	Fullname     string    `json:"fullname"`
	Avatar       string    `json:"avatar"`
	State        int64     `json:"state"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdateAt     time.Time `json:"update_at"`
}