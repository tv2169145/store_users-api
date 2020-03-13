package users

import (
	"errors"
	"fmt"
	"github.com/tv2169145/store_users-api/datasources/mysql/users_db"
	"github.com/tv2169145/store_users-api/logger"
	"github.com/tv2169145/store_users-api/utils/mysql_utils"
	"github.com/tv2169145/store_utils-go/rest_errors"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

const (
	queryInsertUser       = "INSERT INTO users ( first_name, last_name, email, date_created, password, status) VALUES (?, ?, ?, ?, ?, ?);"
	queryGetUser          = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE id=?"
	queryUpdate           = "UPDATE users SET first_name=?, last_name=?, email=? WHERE id=?;"
	queryDelete           = "DELETE FROM users WHERE id=?;"
	queryFindUserByStatus = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE status=?;"
	queryFindByEmailAndPassword = "SELECT id, first_name, last_name, email, password, date_created, status FROM users WHERE email=? AND status=?;"
)

func (u *User) Get() *rest_errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryGetUser)
	if err != nil {
		logger.Error("error when trying to prepare get user statement", err)
		return rest_errors.NewInternalServerError("error when trying to get user", errors.New("database error"))
	}
	defer stmt.Close()

	result := stmt.QueryRow(u.Id)
	if getErr := result.Scan(&u.Id, &u.FirstName, &u.LastName, &u.Email, &u.DateCreated, &u.Status); getErr != nil {
		logger.Error("error when trying to get user by id", err)
		return rest_errors.NewInternalServerError("error when trying to get user", errors.New("database error"))
	}
	return nil
}

func (u *User) Save() *rest_errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryInsertUser)
	if err != nil {
		logger.Error("error when trying to prepare save user statement", err)
		return rest_errors.NewInternalServerError("error when trying to save user", errors.New("database error"))
	}
	defer stmt.Close()
	//result, err := users_db.Client.Exec(queryInsertUser, u.FirstName, u.LastName, u.Email, date_utils.GetNowString())

	insertResult, saveErr := stmt.Exec(u.FirstName, u.LastName, u.Email, u.DateCreated, u.Password, u.Status)
	if saveErr != nil {
		logger.Error("error when trying to save user", err)
		return rest_errors.NewInternalServerError("error when trying to save user", errors.New("database error"))
	}
	userId, err := insertResult.LastInsertId()
	if err != nil {
		logger.Error("error when trying to get last insert id after creating a new user", err)
		return rest_errors.NewInternalServerError("error when trying to save user", errors.New("database error"))
	}
	u.Id = userId
	return nil
}

func (u *User) Update() *rest_errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryUpdate)
	if err != nil {
		logger.Error("error when trying to prepare update user statement", err)
		return rest_errors.NewInternalServerError("error when trying to update user", errors.New("database error"))
	}
	defer stmt.Close()
	_, err = stmt.Exec(u.FirstName, u.LastName, u.Email, u.Id)
	if err != nil {
		logger.Error("error when trying to update user", err)
		return rest_errors.NewInternalServerError("error when trying to update user", errors.New("database error"))
	}
	return nil
}

func (u *User) Delete() *rest_errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryDelete)
	if err != nil {
		logger.Error("error when trying to prepare delete user statement", err)
		return rest_errors.NewInternalServerError("error when trying to delete user", errors.New("database error"))
	}
	defer stmt.Close()
	if _, err = stmt.Exec(u.Id); err != nil {
		logger.Error("error when trying to delete user", err)
		return rest_errors.NewInternalServerError("error when trying to delete user", errors.New("database error"))
	}
	return nil
}

func (u *User) FindByStatus(status string) (Users, *rest_errors.RestErr) {
	stmt, err := users_db.Client.Prepare(queryFindUserByStatus)
	if err != nil {
		logger.Error("error when trying prepare to get active users", err)
		return nil, rest_errors.NewInternalServerError("error when trying to get user", errors.New("database error"))
	}
	defer stmt.Close()

	rows, err := stmt.Query(status)
	if err != nil {
		logger.Error("error when trying to get active users", err)
		return nil, rest_errors.NewInternalServerError("error when trying to get user", errors.New("database error"))
	}
	defer rows.Close()

	results := make([]User, 0)
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); err != nil {
			logger.Error("error when scan user row into user struct", err)
			return nil, rest_errors.NewInternalServerError("error when trying to get user", errors.New("database error"))
		}
		results = append(results, user)
	}
	if len(results) == 0 {
		return nil, rest_errors.NewNotFoundError(fmt.Sprintf("no users matching status: %s", status))
	}
	return results, nil
}

func (u *User) FindByEmailAndPassword() *rest_errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryFindByEmailAndPassword)
	if err != nil {
		logger.Error("error when trying to prepare get user by email and password statement", err)
		return rest_errors.NewInternalServerError("database error", err)
	}
	defer stmt.Close()
	inputPassword := u.Password
	result := stmt.QueryRow(u.Email, StatusActive)
	if getErr := result.Scan(&u.Id, &u.FirstName, &u.LastName, &u.Email, &u.Password, &u.DateCreated, &u.Status); getErr != nil {
		if strings.Contains(getErr.Error(), mysql_utils.ErrorNoRows) {
			return rest_errors.NewNotFoundError("no user found in db")
		}
		logger.Error("error when trying to get user by email and password", err)
		return rest_errors.NewInternalServerError("database error", err)
	}
	if authErr := u.Authenticate(inputPassword); authErr != nil {
		return authErr
	}
	//u.Id = getUser.Id
	//u.FirstName = getUser.FirstName
	//u.LastName = getUser.LastName
	//u.Email = getUser.Email
	//u.DateCreated = getUser.DateCreated
	//u.Status = getUser.Status
	return nil
}

func (u *User) Authenticate(password string) *rest_errors.RestErr {
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return rest_errors.NewBadRequestError("password error")
	}
	return nil
}
