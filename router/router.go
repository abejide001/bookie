package router

import (
	"bookie/controllers"
	"bookie/driver"
	"bookie/models"
	"database/sql"
	"fmt"
	"net/http"
	
	"github.com/gorilla/mux"
)

var books []models.Book
var db *sql.DB

// Router function
func Router() {
	db = driver.ConnectDB()
	controller := controllers.Controller{}
	req := mux.NewRouter()
	req.HandleFunc("/books", controller.GetBooks(db)).Methods("GET")
	req.HandleFunc("/books/{id}", controller.GetBook(db)).Methods("GET")
	req.HandleFunc("/books", controller.AddBook(db)).Methods("POST")
	req.HandleFunc("/books/{id}", controller.DeleteBook(db)).Methods("DELETE")
	req.HandleFunc("/books/{id}", controller.UpdateBook(db)).Methods("PUT")
	errht := http.ListenAndServe(":4000", req)
	if errht != nil {
		fmt.Println("there is an error with http", errht)
		return
	}
}
