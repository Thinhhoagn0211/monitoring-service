package model

import "time"

type User struct {
	Id           int       `json:"id" orm:"column(id)"`
	Email        string    `json:"email" orm:"column(email)" validate:"required"`
	User         string    `json:"username" orm:"column(username)" validate:"required"`
	Password     string    `json:"password" orm:"-"`
	PasswordHash string    `json:"-" orm:"column(password)"`
	Phone        string    `json:"phone" orm:"column(phone)"`
	FullName     string    `json:"full_name" orm:"column(full_name)"`
	Avatar       string    `json:"avatar" orm:"column(avatar)"`
	State        UserState `json:"state" orm:"column(state)"`
	Role         string    `json:"role" orm:"column(role)"`
	CreatedAt    time.Time `json:"create_at" orm:"column(created_at);auto_now_add"`
	UpdateAt     time.Time `json:"update_at" orm:"column(update_at);auto_now"`
}
type UserState = uint8

const (
	UserStateActive UserState = iota
	UserStateBanned
	UserStateDeleted
)
