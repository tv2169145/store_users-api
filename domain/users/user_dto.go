package users

import (
	"github.com/tv2169145/store_utils-go/rest_errors"
	"strings"
)

const (
	StatusActive = "active"
)

type User struct {
	Id          int64  `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	DateCreated string `json:"date_created"`
	Status      string `json:"status"`
	Password    string `json:"password"`
}

type Users []User

func (u *User) Validate() *rest_errors.RestErr {
	u.FirstName = strings.TrimSpace(u.FirstName)
	u.LastName = strings.TrimSpace(u.LastName)
	u.Email = strings.TrimSpace(strings.ToLower(u.Email))
	u.Password = strings.TrimSpace(u.Password)

	if u.Email == "" {
		return rest_errors.NewBadRequestError("invalid email address")
	}
	if u.Password == "" {
		return rest_errors.NewBadRequestError("invalid password")
	}
	return nil
}
