package main

import (
	"database/sql"
	"fmt"
	_ "github.com/glebarez/sqlite"
)

func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "forum.db")
	if err != nil {
		return nil, err
	}

	tables := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			email TEXT UNIQUE NOT NULL,
			username TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL,
			role INTEGER NOT NULL
		);`,

		`CREATE TABLE IF NOT EXISTS posts (
			id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			status INTEGER NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY(user_id) REFERENCES users(id)
		);`,

		`CREATE TABLE IF NOT EXISTS comments (
			id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL,
			post_id TEXT NOT NULL,
			content TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY(user_id) REFERENCES users(id),
			FOREIGN KEY(post_id) REFERENCES posts(id)
		);`,

		`CREATE TABLE IF NOT EXISTS reports (
			id TEXT PRIMARY KEY,
			post_id TEXT NOT NULL,
			reporter_id TEXT NOT NULL,
			reason TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY(post_id) REFERENCES posts(id),
			FOREIGN KEY(reporter_id) REFERENCES users(id)
		);`,
	}

	for _, table := range tables {
		_, err = db.Exec(table)
		if err != nil {
			return nil, fmt.Errorf("failed to create table: %v", err)
		}
	}

	fmt.Println("Database initialized successfully!")
	return db, nil
} 