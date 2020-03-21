package users_db

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/tv2169145/store_utils-go/logger"
	"log"
	"os"
)

const (
	mysql_users_username = "mysql_users_username"
	mysql_users_password = "mysql_users_password"
	mysql_users_host = "mysql_users_host"
	mysql_users_schema = "mysql_users_schema"
)

var (
	Client *sql.DB

	username = os.Getenv(mysql_users_username)
	password = os.Getenv(mysql_users_password)
	host = os.Getenv(mysql_users_host)
	schema = os.Getenv(mysql_users_schema)
)

func init() {
	//datasourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8", username, password, host, schema)
	var err error
	//Client, err = sql.Open("mysql", datasourceName)
	Client, err = sql.Open("mysql", "root:12345678@tcp(localhost:3306)/users_db")
	if err != nil {
		panic(err)
	}
	if err = Client.Ping(); err != nil {
		fmt.Println("1123123")
		panic(err)
	}
	mysql.SetLogger(logger.GetLogger())
	log.Println("database successfully configured")
}

