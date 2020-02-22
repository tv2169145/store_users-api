package users_test

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/tv2169145/store_users-api/datasources/mysql/users_db"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	err   error
	db    *sql.DB
	dbSql *sql.DB
	mock  sqlmock.Sqlmock

	//gr repos.GlobalRepository

	truncateUsers = func() {
		//mock.ExpectQuery("TRUNCATE users").WillReturnRows(sqlmock.NewRows([]string{}))
		//_, err = dbSql.Query("TRUNCATE users")
		//Ω(err).To(BeNil())
	}

	clearDatabase = func() {
		if dbSql == nil {
			Fail("unable to run test because database is missing")
		}
		truncateUsers()
		return
	}
)

var (
	_ = BeforeSuite(func() {
		// connection string - root:pass@tcp(localhost:3306)/grpc
		// root:12345678@tcp(localhost:3306)/grpc
		//db, err = sql.Open("mysql", "")
		//if err != nil {
		//	panic(err)
		//}
		//Ω(err).To(BeNil())
		dbSql, mock, err = sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		Ω(err).To(BeNil())
		db = dbSql
		users_db.Client = dbSql
	})

	_ = AfterSuite(func() {
		//err = mock.ExpectationsWereMet()
		//Ω(err).To(BeNil())
	})
)

func TestUsers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Users Suite")
}
