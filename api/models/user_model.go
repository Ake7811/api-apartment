package models

import (
	"apartment/api/security"
	"errors"
	"html"
	"strings"
	"time"
)

type User struct {
	Id          uint      `gorm:"primary_key" json:"id"`
	Username    string    `gorm:"size:100;not null" json:"username"`
	Password    string    `gorm:"size:200;not null" json:"password"`
	Firstname   string    `gorm:"size:100;not null" json:"firstname"`
	Lastname    string    `gorm:"size:100" json:"lastname"`
	User_status string    `gorm:"size:1;not null" json:"user_status"`
	User_level  string    `gorm:"size:1;not null" json:"user_level"`
	Create_At   time.Time `gorm:"default:current_timestamp()" json:"create_at"`
	Updated_At  time.Time `gorm:"default:current_timestamp()" json:"updated_at"`
}

func (u *User) BeforeSave() error {
	hashedPassword, err := security.Hash(u.Password)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) Prepare() {
	u.Id = 0
	u.Username = html.EscapeString(strings.TrimSpace(u.Username))
	u.Password = html.EscapeString(strings.TrimSpace(u.Password))
	u.Firstname = html.EscapeString(strings.TrimSpace(u.Firstname))
	u.Lastname = html.EscapeString(strings.TrimSpace(u.Lastname))
	u.Create_At = time.Now()
	u.Updated_At = time.Now()
}

func (u *User) Validate(action string) error {

	switch strings.ToLower(action) {
	case "update":
		if u.Username == "" {
			return errors.New("username Required")
		}
		if u.Password == "" {
			return errors.New("password Required")
		}
		if u.Firstname == "" {
			return errors.New("firstname Required")
		}
		return nil
	case "login":
		if u.Username == "" {
			return errors.New("username Required")
		}
		if u.Password == "" {
			return errors.New("password Required")
		}
		return nil
	default:
		if u.Username == "" {
			return errors.New("username Required")
		}
		if u.Password == "" {
			return errors.New("password Required")
		}
		if u.Firstname == "" {
			return errors.New("firstname Required")
		}
		return nil

	}

}
