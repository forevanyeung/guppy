package cmd

import (
	"fmt"

	"github.com/forevanyeung/guppy/cf"
	"github.com/spf13/cobra"
	"google.golang.org/api/drive/v3"
)

func init() {
	rootCmd.AddCommand(uploadCmd)
}

var uploadCmd = &cobra.Command{
	Use:   "upload [file]",
	Short: "Upload a file to Google Drive",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filePath := args[0]
		upload(filePath)
	},
}

var listenPort int = 8080

type GuppyStatus struct {
	IsAuthenticated bool   `json:"isAuthenticated"`
	UploadStarted   bool   `json:"uploadStarted"`
	UploadFinished  bool   `json:"uploadFinished"`
	UploadError     string `json:"uploadError"`
	WebLink         string `json:"webLink"`
}

var driveService *drive.Service

func upload(filePath string) {
	// TODO pick a random open port
	listenPort = 8080

	// Get the Google OAuth2 client ID from the configuration
	domain := "com.forevanyeung.guppy"
	value := cf.CFPreferencesCopyAppValue("GoogleOauth2ClientId", domain)
	if value == nil {
		fmt.Printf("Missing %s configuration\n", domain)
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

	// Start the upload loop with the file path as an argument
	uploadFile(filePath, statusChan)

	// Wait for the HTTP server to finish before exiting
	<-serverDone	
}
