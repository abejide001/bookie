package driver

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"log"
	"os"
)

var db *sql.DB

// ConnectDB function
func ConnectDB() *sql.DB {
	pgURL, pgerr := pq.ParseURL(os.Getenv("ELEPHANT_SQL_URL"))
	if pgerr != nil {
		log.Fatal(pgerr)
	}
	db, err := sql.Open("postgres", pgURL)
	if err != nil {
		fmt.Println("db is not connected")
	}
	defer db.Close()
	dberr := db.Ping()
	if dberr != nil {
		fmt.Println(dberr)
	}
	return db
}
