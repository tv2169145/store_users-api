package users

import (
	"fmt"
	"github.com/tv2169145/store_users-api/utils/errors"
)

var (
	userDB = make(map[int64]*User)
)

func (u *User) Get() *errors.RestErr {
	result := userDB[u.Id]
	if result == nil {
		return errors.NewNotFoundError(fmt.Sprintf("user: %d not found", u.Id))
	}
	u.Id = result.Id
	u.Email = result.Email
	u.FirstName = result.FirstName
	u.LastName = result.LastName
	u.DateCreated = result.DateCreated
	return nil
}

func (u *User) Save() *errors.RestErr {
	current := userDB[u.Id]
	if current != nil {
		if current.Email == u.Email {
			return errors.NewBadRequestError(fmt.Sprintf("email: %s already registered", u.Email))
		}
		return errors.NewBadRequestError(fmt.Sprintf("user: %d is readly exist", u.Id))
	}
	userDB[u.Id] = u
	return nil
}
