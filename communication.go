package main

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

func createPost(w http.ResponseWriter, r *http.Request) {
	sessionCookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	mu.Lock()
	session, exists := sessions[sessionCookie.Value]
	mu.Unlock()
	if !exists {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID := session.UserID
	title := r.FormValue("title")
	content := r.FormValue("content")
	categories := r.Form["categories"]

	// Handle image upload
	imageURL := ""
	file, handler, err := r.FormFile("image")
	if err == nil {
		defer file.Close()

		ext := strings.ToLower(filepath.Ext(handler.Filename))
		if ext == ".jpeg" || ext == ".jpg" || ext == ".png" || ext == ".gif" {
			fileName := uuid.New().String() + ext
			filePath := filepath.Join("uploads", fileName)

			err = os.MkdirAll("uploads", os.ModePerm)
			if err == nil {
				out, err := os.Create(filePath)
				if err == nil {
					defer out.Close()
					_, err = io.Copy(out, file)
					if err == nil {
						imageURL = "/uploads/" + fileName
					}
				}
			}
		}
	}

	post := Post{
		ID:         uuid.New().String(),
		UserID:     userID,
		Title:      title,
		Content:    content,
		Categories: categories,
		ImageURL:   imageURL,
		Status:     Pending, // Default status for new posts
	}

	mu.Lock()
	posts[post.ID] = post
	mu.Unlock()

	w.Write([]byte("Post created successfully"))
}

func createComment(w http.ResponseWriter, r *http.Request) {
	sessionCookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	mu.Lock()
	session, exists := sessions[sessionCookie.Value]
	mu.Unlock()
	if !exists {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID := session.UserID
	postID := r.FormValue("post_id")
	content := r.FormValue("content")

	comment := Comment{
		ID:      uuid.New().String(),
		UserID:  userID,
		PostID:  postID,
		Content: content,
	}

	mu.Lock()
	comments[comment.ID] = comment
	mu.Unlock()

	w.Write([]byte("Comment created successfully"))
}
