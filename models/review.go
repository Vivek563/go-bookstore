package models

import (
	"context"
	"database/sql"
	"github.com/redis/go-redis/v9"
	"strconv"
)

type Review struct {
	Isbn         string  `json:"isbn"`
	Title        string  `json:"title"`
	Rating       float32 `json:"rating"`
	Price        float32 `json:"price"`
	Descriptions string  `json:"descriptions"`
	Uid          int64   `json:"uid"`
}

type ReviewModel struct {
	db     *sql.DB
	redis  *redis.Client
	ctx    context.Context
	cancel context.CancelFunc
}

// create new review
func (m *ReviewModel) CreateReview(review *Review) error {
	_, err := m.db.Exec("INSERT INTO `reviews` (`isbn`,`price`,`title`,`rating`,`descriptions`,`uid`) VALUES (?,?,?,?,?,? )", &review.Isbn, &review.Price, &review.Title, &review.Rating, &review.Descriptions, &review.Uid)
	return err
}

func (m *ReviewModel) DeleteReview(id, uid int64) error {
	_, err := m.db.Exec("UPDATE `reviews` SET is_deleted = 1 WHERE  `id` = ? AND uid = ? ", id, uid)
	return err
}

func (m *ReviewModel) ReviewListing() ([]*Review, error) {
	stmt := "SELECT isbn, title, rating, price, descriptions, uid FROM `reviews` WHERE is_deleted = 0"
	reviews, err := m.Listing(stmt)
	return reviews, err
}

// if user logged it will show its review
func (m *ReviewModel) MyReview(uid int64) ([]*Review, error) {
	stmt := "SELECT isbn, title, rating, price, descriptions, uid FROM reviews WHERE is_deleted = 0 AND  uid = '" + strconv.Itoa(int(uid)) + "'"
	reviews, err := m.Listing(stmt)
	return reviews, err
}

func (m *ReviewModel) Listing(stmt string) ([]*Review, error) {
	rows, err := m.db.Query(stmt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}

	defer rows.Close()

	reviews := []*Review{}
	for rows.Next() {
		bk, err := m.ScanReviewData(rows)
		if err != nil {
			return nil, err
		}

		reviews = append(reviews, bk)
	}

	return reviews, err
}

func (m *ReviewModel) GetReviewByIsbn(isbn string) ([]*Review, error) {
	stmt := "SELECT isbn, title, rating, price, descriptions, uid FROM reviews WHERE isbn = '" + isbn + "' AND is_deleted = 0"
	review, err := m.Listing(stmt)
	return review, err
}

func (m *ReviewModel) ScanReviewData(rows *sql.Rows) (*Review, error) {
	review := new(Review)
	err := rows.Scan(
		&review.Isbn,
		&review.Title,
		&review.Rating,
		&review.Price,
		&review.Descriptions,
		&review.Uid)
	return review, err
}
