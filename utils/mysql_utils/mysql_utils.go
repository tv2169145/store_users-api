package mysql_utils

import (
	"errors"
	"github.com/go-sql-driver/mysql"
	"github.com/tv2169145/store_utils-go/rest_errors"
	"strings"
)

const (
	ErrorNoRows = "no rows in result set"
)

func ParseError(err error) *rest_errors.RestErr {
	sqlErr, ok := err.(*mysql.MySQLError)
	if !ok {
		if strings.Contains(err.Error(), ErrorNoRows) {
			return rest_errors.NewNotFoundError("no record match given id")
		}
		return rest_errors.NewInternalServerError("error parsing database response", err)
	}
	switch sqlErr.Number {
	case 1062:
		return rest_errors.NewBadRequestError("invalid data")
	}
	return rest_errors.NewInternalServerError("error processing request", errors.New("database error"))
}
