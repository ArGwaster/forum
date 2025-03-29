package main

import (
	"net/http"

	"github.com/google/uuid"
)

func approvePost(w http.ResponseWriter, r *http.Request) {
	postID := r.FormValue("post_id")

	mu.Lock()
	defer mu.Unlock()

	post, exists := posts[postID]
	if !exists {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	post.Status = Approved
	posts[postID] = post

	w.Write([]byte("Post approved successfully"))
}

func rejectPost(w http.ResponseWriter, r *http.Request) {
	postID := r.FormValue("post_id")

	mu.Lock()
	defer mu.Unlock()

	post, exists := posts[postID]
	if !exists {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	post.Status = Rejected
	posts[postID] = post

	w.Write([]byte("Post rejected successfully"))
}

func reportPost(w http.ResponseWriter, r *http.Request) {
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

	reporterID := session.UserID
	postID := r.FormValue("post_id")
	reason := r.FormValue("reason")

	report := Report{
		ID:         uuid.New().String(),
		PostID:     postID,
		ReporterID: reporterID,
		Reason:     reason,
	}

	mu.Lock()
	reports[report.ID] = report
	mu.Unlock()

	w.Write([]byte("Post reported successfully"))
}
