package handlers

import (
	"log"
	"net/http"
)

func HandleInternalError(w http.ResponseWriter, err error) {
	// Log the internal error
	log.Printf("Internal error: %v", err)

	// Return a generic error message to the client
	http.Error(w, "An unexpected error occurred. Please try again later.", http.StatusInternalServerError)
}
