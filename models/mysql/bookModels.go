package mysql

import (
	"github.com/afifallimaruf/go-booksell/config"
	"github.com/afifallimaruf/go-booksell/models"
)

// model untuk menyimpan data buku ke database
func InsertBook(title, author, summary, price, imgName, category string) bool {
	query := "INSERT INTO books (title, author, summary, price, image, category) VALUES (?, ?, ?, ?, ?, ?)"

	db, err := config.ConnectDB()
	if err != nil {
		panic(err)
	}

	result, err := db.Exec(query, title, author, summary, price, imgName, category)
	if err != nil {
		panic(err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return false
	}

	return id > 0
}

func GetBooks() ([]*models.Books, error) {
	query := "SELECT * FROM books"

	db, err := config.ConnectDB()
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	books := []*models.Books{}

	for rows.Next() {
		book := &models.Books{}

		rows.Scan(&book.Id, &book.Title, &book.Author, &book.Summary, &book.Price, &book.Image, &book.Category)

		books = append(books, book)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return books, nil
}

// Ambil data buku bergenre non fiksi
func GetNonFiction() (*models.Books, error) {
	db, err := config.ConnectDB()
	if err != nil {
		return nil, err
	}

	query := `SELECT * FROM books WHERE category = "Non-Fiction"`

	row := db.QueryRow(query)

	book := &models.Books{}

	if err := row.Scan(&book.Id, &book.Title, &book.Author, &book.Summary, &book.Price, &book.Image, &book.Category); err != nil {
		return nil, err
	}

	return book, nil
}
