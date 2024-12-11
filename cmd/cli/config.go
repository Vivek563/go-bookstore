package main

import (
	"database/sql" // db
)

// for a given DSN
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

// func Init() error {
// 	return createAccountTable()
// }

// func createAccountTable(s *sql.DB) error {
// 	query := `create table if not exists bookreviewss (
// 		id INT PRIMARY KEY AUTOINCREMENT,
// 		isbn VARCHAR(50) NOT NULL,
// 		title VARCHAR(100) NOT NULL,
// 		author VARCHAR(50) NOT NULL,
// 		descriptions TEXT NOT NULL,
// 		uid INT NOT NULL,
// 		created_at DATETIME DEFAULT CURRENTTIMESTAMP
// 	);`

// 	_, err := s.db.Exec(query)
// 	return err
// }
