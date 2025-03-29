package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

const maxImageSize = 20 * 1024 * 1024 // 20 MB

func uploadImage(w http.ResponseWriter, r *http.Request) {
	// Parse the multipart form
	err := r.ParseMultipartForm(maxImageSize)
	if err != nil {
		http.Error(w, "Image too big", http.StatusBadRequest)
		return
	}

	// Get the file from the form
	file, handler, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Failed to upload image", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Validate the file type
	ext := strings.ToLower(filepath.Ext(handler.Filename))
	if ext != ".jpeg" && ext != ".jpg" && ext != ".png" && ext != ".gif" {
		http.Error(w, "Invalid file type", http.StatusBadRequest)
		return
	}

	// Generate a unique file name
	fileName := uuid.New().String() + ext
	filePath := filepath.Join("uploads", fileName)

	// Create the uploads directory if it doesn't exist
	err = os.MkdirAll("uploads", os.ModePerm)
	if err != nil {
		http.Error(w, "Failed to create uploads directory", http.StatusInternalServerError)
		return
	}

	// Save the file
	out, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Failed to save image", http.StatusInternalServerError)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		http.Error(w, "Failed to save image", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(fmt.Sprintf("Image uploaded successfully: %s", fileName)))
}
