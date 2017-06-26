package main

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"log"
	"net/http"
	"os"
)

func main() {
	// Read password from enviroment
	const env = "VOID_FILE_EXPLORER_PASSWORD"
	password := os.Getenv(env)

	// Generate password if it wasn't specified by environment variable
	if password == "" {
		log.Println(env, "not set. Generating password...")
		var secret [18]byte
		if _, err := rand.Read(secret[:]); err != nil {
			log.Fatalln("Could not generate password:", err.Error())
		}
		password = base64.URLEncoding.EncodeToString(secret[:])
	}

	// Write login details to log
	log.Println("Any user is valid.")
	log.Println("The password is:", password)

	// Handle requests
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Authenticate
		w.Header().Set("WWW-Authenticate", `Basic realm="sourcevoid_diskexplorer"`)
		if _, requestPassword, _ := r.BasicAuth(); subtle.ConstantTimeCompare([]byte(password), []byte(requestPassword)) != 1 {
			http.Error(w, "Enter any user and the password.", http.StatusUnauthorized)
		}
		// Handle file uploads
		if err := Upload(r, w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		// Serve directory contents
		if err := Browser(w, r); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	// Start server
	log.Fatalln(http.ListenAndServe(":8080", nil))
}
