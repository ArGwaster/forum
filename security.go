package main

import (
	"net/http"
	"sync"
	"time"

	"golang.org/x/crypto/acme/autocert"
	"golang.org/x/crypto/bcrypt"
)

// Session struct to store session information

// HashPassword hashes the given password using bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash checks if the given password matches the hashed password
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Rate limiter using a token bucket algorithm
type RateLimiter struct {
	tokens     int
	maxTokens  int
	refillRate time.Duration
	lastRefill time.Time
	mu         sync.Mutex
}

func NewRateLimiter(maxTokens int, refillRate time.Duration) *RateLimiter {
	return &RateLimiter{
		tokens:     maxTokens,
		maxTokens:  maxTokens,
		refillRate: refillRate,
		lastRefill: time.Now(),
	}
}

func (rl *RateLimiter) Allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(rl.lastRefill)
	rl.tokens += int(elapsed / rl.refillRate)
	if rl.tokens > rl.maxTokens {
		rl.tokens = rl.maxTokens
	}
	rl.lastRefill = now

	if rl.tokens > 0 {
		rl.tokens--
		return true
	}
	return false
}

func security() {
	// Setup HTTPS server with autocert
	m := &autocert.Manager{
		Cache:      autocert.DirCache("certs"),
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist("yourdomain.com"),
	}

	server := &http.Server{
		Addr:      ":https",
		TLSConfig: m.TLSConfig(),
		Handler:   setupRoutes(),
	}

	go http.ListenAndServe(":http", m.HTTPHandler(nil))
	server.ListenAndServeTLS("", "")
}

func setupRoutes() http.Handler {
	mux := http.NewServeMux()

	// Serve static files from the "templates" directory
	fs := http.FileServer(http.Dir("./templates"))
	mux.Handle("/templates/", http.StripPrefix("/templates/", fs))

	// Handle the root route to redirect to the register.html page
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/templates/register.html", http.StatusSeeOther)
	})

	// Add more routes here...

	return mux
}
