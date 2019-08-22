// entry point into our application
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/lib/pq"
	"github.com/subosito/gotenv"
	"log"
	"net/http"
	"os"
	"strconv"
)

// Book type
type Book struct {
	ID     int    `json:id`
	Title  string `json:title`
	Author string `json:author`
	Year   int    `json:year`
}

var books []Book
var db *sql.DB

func init() {
	gotenv.Load() // loads env variables
}

func main() {
	pgURL, err := pq.ParseURL(os.Getenv("ELEPHANT_SQL_URL"))
	if err != nil {
		log.Fatal(err)
	}
	db, err = sql.Open("postgres", pgURL)
	if err != nil {
		fmt.Println("db is not connected")
	}
	defer db.Close()
	dberr := db.Ping()
	if dberr != nil {
		fmt.Println(dberr)
	}
	req := mux.NewRouter()
	req.HandleFunc("/books", getBooks).Methods("GET")
	req.HandleFunc("/books/{id}", getBook).Methods("GET")
	req.HandleFunc("/books", addBook).Methods("POST")
	req.HandleFunc("/books/{id}", deleteBook).Methods("DELETE")
	req.HandleFunc("/books/{id}", updateBook).Methods("PUT")
	errht := http.ListenAndServe(":4000", req)
	if errht != nil {
		fmt.Println("there is an error with http", errht)
		return
	}
}

func getBooks(res http.ResponseWriter, req *http.Request) {
	var book Book
	fmt.Println(book)
	books = []Book{}
	rows, err := db.Query("select * from bookstore")
	if err != nil {
		log.Fatal("db error", err)
	}
	for rows.Next() {
		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year)
		if err != nil {
			fmt.Println(err)
			return
		}
		books = append(books, book)
	}
	defer rows.Close()
	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(books)
}

func getBook(res http.ResponseWriter, req *http.Request) {
	var book Book

	params := mux.Vars(req)
	id, _ := strconv.Atoi(params["id"])
	rows := db.QueryRow("select * from bookstore where id=$1", id)
	err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year)
	if id != book.ID {
		var errMessage = map[string]string{"status": "ID does not exist", "code": "400"}
		res.Header().Set("Content-Type", "application/json")

		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(errMessage)
		return
	}
	if err != nil {
		fmt.Println(err)
	}
	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(book)
}

func addBook(res http.ResponseWriter, req *http.Request) {
	var book Book
	var bookID int
	json.NewDecoder(req.Body).Decode(&book) // decode maps the value in the body to the book var
	if len(book.Title) <= 0 {
		var errMessage = map[string]string{"status": "enter a value for title", "code": "400"}
		res.Header().Set("Content-Type", "application/json")

		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(errMessage)
		return
	}
	err := db.QueryRow(
		"insert into bookstore (title, author, year) values($1, $2, $3) RETURNING id;",
		book.Title, book.Author, book.Year,
	).Scan(&bookID)
	if err != nil {
		fmt.Println(err)
	}
	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(book)
}

func updateBook(res http.ResponseWriter, req *http.Request) {
	var book Book
	json.NewDecoder(req.Body).Decode(&book)
	params := mux.Vars(req)
	id, _ := strconv.Atoi(params["id"])
	results, err := db.Exec(
		"update bookstore set title=$1, author=$2, year=$3 where id=$4 RETURNING id",
		&book.Title, &book.Author, &book.Year, &id,
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	rowsUpdated, _ := results.RowsAffected()
	if rowsUpdated == 0 {
		var errMessage = map[string]string{"status": "ID does not exist", "code": "400"}
		res.Header().Set("Content-Type", "application/json")

		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(errMessage)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(book)
}

func deleteBook(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	id, _ := strconv.Atoi(params["id"])
	results, err := db.Exec("delete from bookstore where id=$1", id)
	if err != nil {
		fmt.Println(err)
	}
	rowsDeleted, _ := results.RowsAffected()
	if rowsDeleted == 0 {
		var errMessage = map[string]string{"status": "ID does not exist", "code": "400"}
		res.Header().Set("Content-Type", "application/json")

		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(errMessage)
		return
	}
	json.NewEncoder(res).Encode(rowsDeleted)
}
