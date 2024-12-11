package models

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Book struct {
	ISBN         string  `json:"isbn"`
	Title        string  `json:"title"`
	Author       string  `json:"author"`
	Price        float32 `json:"price"`
	Descriptions string  `json:"descriptions"`
	Genre        string  `json:"genre"`
}

func (m *BookModel) Close() {
	m.cancel()
	m.redis.Close()
	m.db.Close()
}

type BookModel struct {
	db     *sql.DB
	redis  *redis.Client
	ctx    context.Context
	cancel context.CancelFunc
}

// add books in db
func (m *BookModel) CreateBook(book *Book) error {
	_, err := m.db.Exec("INSERT INTO `books` (`isbn`,`title`,`author`,`price`,`descriptions`,`genre`) VALUES (?,?,?,?,?,? )", &book.ISBN, &book.Title, &book.Author, &book.Price, &book.Descriptions, &book.Genre)
	return err
}

// check isbn already exist or not()
func (m *BookModel) BookExist(ISBN string) bool {
	var valid int
	_ = m.db.QueryRow("SELECT 1 FROM `books` WHERE  `isbn` = ?", ISBN).Scan(&valid)
	return valid > 0
}

func (m *BookModel) GetBookByIsbn(ISBN string) ([]*Book, error) {
	stmt := fmt.Sprintf("SELECT `isbn`,`title`,`author`,`price`,`descriptions`,`genre` FROM `books` WHERE `ISBN` = '%s' ", ISBN)
	bk, err := m.Listing(stmt)
	return bk, err
}

func (m *BookModel) BooksListing() ([]*Book, error) {
	stmt := "SELECT `isbn`,`title`,`author`,`price`,`descriptions`,`genre` FROM books"
	Books, err := m.Listing(stmt)
	return Books, err
}

func (m *BookModel) Listing(stmt string) ([]*Book, error) {

	Books := []*Book{}
	queryBytes, err := json.Marshal(stmt)
	if err != nil {
		panic(err)
	}

	val, err := m.redis.Get(m.ctx, string(queryBytes)).Result()
	if err == nil {
		// Deserialize the cached result
		err = json.Unmarshal([]byte(val), &Books)
		if err != nil {
			return nil, err
		}

		return Books, err
	} else if err != redis.Nil {
		return nil, err
	}

	rows, err := m.db.Query(stmt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}

	defer rows.Close()

	for rows.Next() {
		book := new(Book)
		err := m.ScanBookData(rows, book)
		if err != nil {
			return nil, err
		}

		Books = append(Books, book)
	}

	// Cache the result in Redis for 5 minutes
	data, err := json.Marshal(Books)
	if err != nil {
		return nil, err
	}

	err = m.redis.Set(m.ctx, string(queryBytes), data, 5*time.Minute).Err()
	if err != nil {
		return nil, err
	}

	m.Close()
	return Books, err
}

func (m *BookModel) ScanBookData(rows *sql.Rows, book *Book) error {
	return rows.Scan(
		&book.ISBN,
		&book.Title,
		&book.Author,
		&book.Price,
		&book.Descriptions,
		&book.Genre,
	)
}
