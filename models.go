package main

import (
	"sync"
	"time"
)

// UserRole defines the role of a user in the forum
type UserRole int

const (
	Guest UserRole = iota
	Users
	Moderator
	Administrator
)

// User struct to store user information
type User struct {
	ID       string   `gorm:"primaryKey"`
	Email    string   `gorm:"unique;not null"`
	Username string   `gorm:"unique;not null"`
	Password string   `gorm:"not null"` // Hashed password
	Role     UserRole `gorm:"not null"`
}

// PostStatus defines the status of a post
type PostStatus int

const (
	Pending PostStatus = iota
	Approved
	Rejected
)

// Post struct to store post information
type Post struct {
	ID         string     `gorm:"primaryKey"`
	UserID     string     `gorm:"not null"`
	Title      string     `gorm:"not null"`
	Content    string     `gorm:"not null"`
	Categories []string   `gorm:"-"`
	ImageURL   string     `gorm:"-"`
	Status     PostStatus `gorm:"not null"`
	CreatedAt  time.Time  `gorm:"autoCreateTime"`
}

// Comment struct to store comment information
type Comment struct {
	ID        string    `gorm:"primaryKey"`
	UserID    string    `gorm:"not null"`
	PostID    string    `gorm:"not null"`
	Content   string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

// Session struct to store session information
type Session struct {
	UserID    string
	ExpiresAt time.Time
}

// Report struct to store report information
type Report struct {
	ID         string    `gorm:"primaryKey"`
	PostID     string    `gorm:"not null"`
	ReporterID string    `gorm:"not null"`
	Reason     string    `gorm:"not null"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
}
