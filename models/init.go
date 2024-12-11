package models

import (
	"context"
	"database/sql"

	"github.com/redis/go-redis/v9"
)

type Init struct {
	Books  BookModel
	Users  UserModel
	Review ReviewModel
}

func Constructor(db *sql.DB, rd *redis.Client) *Init {
	ctx, cancel := context.WithCancel(context.Background())
	return &Init{
		Books:  BookModel{db: db, redis: rd, ctx: ctx, cancel: cancel},
		Users:  UserModel{db: db, redis: rd, ctx: ctx, cancel: cancel},
		Review: ReviewModel{db: db, redis: rd, ctx: ctx, cancel: cancel},
	}
}
