package main

import (
	"fmt"
	"os"

	"google.golang.org/api/drive/v3"
)

var listenPort int = 8080

type GuppyStatus struct {
	IsAuthenticated bool   `json:"isAuthenticated"`
	UploadStarted   bool   `json:"uploadStarted"`
	UploadFinished  bool   `json:"uploadFinished"`
	UploadError     string `json:"uploadError"`
	WebLink         string `json:"webLink"`
}

var driveService *drive.Service

func main() {
	// TODO pick a random open port
	listenPort = 8080

	// Get the Google OAuth2 client ID from the configuration
	domain := "com.forevanyeung.guppy"
	value := CFPreferencesCopyAppValue("GoogleOauth2ClientId", domain)
	if(value == nil) {
		fmt.Printf("Missing %s configuration\n", domain)
		return
	}
	
	// Check if a file path is provided as a command line argument
	if len(os.Args) < 2 {
		fmt.Println("Please provide a file path as a command line argument")
		return
	}

	// Initialize the channel before starting the HTTP server
	newAuthChan := make(chan string)
	statusChan := make(chan GuppyStatus)
	serverDone := make(chan bool)

	// Start the HTTP server
	go httpServer(newAuthChan, statusChan, serverDone)

	// Get auth from keyring or start oauth flow
	auth(newAuthChan, value.(string))

	// Get the file path from the first command line argument
	filePath := os.Args[1]

	// Start the upload loop with the file path as an argument
	uploadFile(filePath, statusChan)

	// Wait for the HTTP server to finish before exiting
	<-serverDone
}
