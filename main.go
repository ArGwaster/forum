package main

import (
	"net/http"
)

func main() {
	// Serve static files from the "templates" directory
	fs := http.FileServer(http.Dir("./templates"))
	http.Handle("/templates/", http.StripPrefix("/templates/", fs))

	// Serve static files from the "uploads" directory
	fsUploads := http.FileServer(http.Dir("./uploads"))
	http.Handle("/uploads/", http.StripPrefix("/uploads/", fsUploads))

	// Handle the root route to redirect to the register.html page
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/templates/home.html", http.StatusSeeOther)
	})

	// Authentication routes
	http.HandleFunc("/register", register)
	http.HandleFunc("/login", login)

	// Communication routes
	http.HandleFunc("/createPost", createPost)
	http.HandleFunc("/createComment", createComment)

	// Image upload route
	http.HandleFunc("/upload", uploadImage)

	// Moderation routes
	http.HandleFunc("/approvePost", approvePost)
	http.HandleFunc("/rejectPost", rejectPost)
	http.HandleFunc("/reportPost", reportPost)

	// Start the server
	http.ListenAndServe(":999", nil)
}
