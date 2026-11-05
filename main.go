package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

var books = []Book{
	{ID: 1, Title: "Golang", Author: "Ali"},
	{ID: 2, Title: "Python", Author: "Mehdi"},
	{ID: 3, Title: "Java", Author: "Hosein"},
}

var nextID = 4

func main() {
	http.HandleFunc("/books", handleBooks)
	http.HandleFunc("/books/search", searchBooks)
	http.HandleFunc("/books/", deleteBook)

	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}

// نمایش همه کتاب ها (GET) و افزودن کتاب (POST)
func handleBooks(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		json.NewEncoder(w).Encode(books)
		return
	}

	if r.Method == "POST" {
		var newBook Book

		err := json.NewDecoder(r.Body).Decode(&newBook)
		if err != nil {
			http.Error(w, "Invalid data", 400)
			return
		}

		newBook.ID = nextID
		nextID = nextID + 1

		books = append(books, newBook)

		json.NewEncoder(w).Encode(newBook)
		return
	}

	http.Error(w, "Method not allowed", 405)
}

// جستجو: /books/search?q=go
func searchBooks(w http.ResponseWriter, r *http.Request) {
	word := strings.ToLower(r.URL.Query().Get("q"))

	var result []Book

	for _, b := range books {
		title := strings.ToLower(b.Title)
		author := strings.ToLower(b.Author)

		if strings.Contains(title, word) || strings.Contains(author, word) {
			result = append(result, b)
		}
	}

	if len(result) == 0 {
		http.Error(w, "Book not found", 404)
		return
	}

	json.NewEncoder(w).Encode(result)
}

// حذف کتاب: DELETE /books/2
func deleteBook(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	idText := strings.TrimPrefix(r.URL.Path, "/books/")

	id, err := strconv.Atoi(idText)
	if err != nil {
		http.Error(w, "Invalid ID", 400)
		return
	}

	for i, b := range books {
		if b.ID == id {
			books = append(books[:i], books[i+1:]...)
			fmt.Fprintln(w, "Book deleted successfully")
			return
		}
	}

	http.Error(w, "Book not found", 404)
}
