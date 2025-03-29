package main

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func hashPassword(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return hex.EncodeToString(hash.Sum(nil))
}

func register(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	username := r.FormValue("username")
	password := r.FormValue("password")

	mu.Lock()
	defer mu.Unlock()

	if _, exists := users[email]; exists {
		http.Error(w, "Email already taken", http.StatusConflict)
		return
	}

	hashedPassword := hashPassword(password)
	userID := uuid.New().String()
	users[email] = User{
		ID:       userID,
		Email:    email,
		Username: username,
		Password: hashedPassword,
		Role:     Users, // Default role for new users
	}
	w.Write([]byte("User registered successfully"))
}

func login(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	mu.Lock()
	defer mu.Unlock()

	user, exists := users[email]
	if !exists || user.Password != hashPassword(password) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	sessionID := uuid.New().String()
	session := Session{
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	sessions[sessionID] = session

	http.SetCookie(w, &http.Cookie{
		Name:    "session_id",
		Value:   sessionID,
		Expires: session.ExpiresAt,
	})

	w.Write([]byte("Login successful"))
}
