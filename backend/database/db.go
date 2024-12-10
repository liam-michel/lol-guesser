package database

import (
	"database/sql"
	"fmt"
	"log"
	_ "github.com/go-sql-driver/mysql"

	
)

var db *sql.DB

func InitDB() {
	var err error
	dbPassword = os.Getenv("MYSQL_ROOT_PASSWORD")
}