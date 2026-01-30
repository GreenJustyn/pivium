package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

// Get secret from environment variable
var secretToken = os.Getenv("PIVIUM_UPDATER_SECRET")

func updateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Simple security: check for a secret in a header
	if secretToken != "" {
		payload, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading request body: %v", err)
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
			return
		}

		// GitHub sends the signature in X-Hub-Signature-256 header
		signature := r.Header.Get("X-Hub-Signature-256")
		if !isValidSignature(signature, payload) {
			log.Println("Invalid signature")
			http.Error(w, "Invalid signature", http.StatusUnauthorized)
			return
		}
	}


	log.Println("Received valid update request. Running update script...")
	// Run the update script in the background
	cmd := exec.Command("/bin/bash", "/opt/pivium/scripts/update.sh")
	cmd.Stdout = log.Writer()
	cmd.Stderr = log.Writer()

	err := cmd.Start()
	if err != nil {
		log.Printf("Error starting update script: %v", err)
		http.Error(w, "Failed to start update", http.StatusInternalServerError)
		return
	}

	// Don't wait for the script to finish. The service will be restarted.
	go func() {
		err := cmd.Wait()
		if err != nil {
			log.Printf("Update script finished with error: %v", err)
		} else {
			log.Println("Update script finished successfully.")
		}
	}()

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("Update process started in background.\n"))
}

func isValidSignature(signatureHeader string, payload []byte) bool {
	if signatureHeader == "" {
		return false
	}
	parts := strings.SplitN(signatureHeader, "=", 2)
	if len(parts) != 2 || parts[0] != "sha256" {
		return false
	}
	
	mac := hmac.New(sha256.New, []byte(secretToken))
	mac.Write(payload)
	expectedMAC := hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(parts[1]), []byte(expectedMAC))
}


func main() {
	if secretToken == "" {
		log.Println("Warning: PIVIUM_UPDATER_SECRET is not set. The updater endpoint will not be secured.")
	}
	http.HandleFunc("/update", updateHandler)
	log.Println("Starting pivium-updater on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
