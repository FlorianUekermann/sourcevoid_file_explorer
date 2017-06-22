package main

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	// Read password from enviroment
	const env = "VOID_DISKEXPLORER_PASSWORD"
	password := os.Getenv(env)

	// Generate password if it wasn't specified by environment variable
	if password == "" {
		fmt.Println(env, "not set. Generating password...")
		var secret [18]byte
		if _, err := rand.Read(secret[:]); err != nil {
			log.Fatalln("Could not generate password:", err.Error())
		}
		password = base64.URLEncoding.EncodeToString(secret[:])
	}

	// Write login details to log
	fmt.Println("Any user is valid.")
	fmt.Println("The password is:", password)

	// Authenticate and serve filesystem
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("WWW-Authenticate", `Basic realm="sourcevoid_diskexplorer"`)
		_, requestPassword, _ := r.BasicAuth()
		if subtle.ConstantTimeCompare([]byte(password), []byte(requestPassword)) != 1 {
			http.Error(w, "Enter any user and the password.", http.StatusUnauthorized)
		} else {
			http.FileServer(http.Dir("/home/cuser")).ServeHTTP(w, r)
		}
	})

	// Start server
	log.Fatalln(http.ListenAndServe(":8080", nil))
}
