package cmd

import (
	"log/slog"
	"net"

	"github.com/forevanyeung/guppy/cli/analytics"
	"github.com/forevanyeung/guppy/cli/cf"
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
		analytics.TrackEvent("$pageview", map[string]interface{}{
			"$current_url": "upload",
		})
		
		filePath := args[0]
		upload(filePath)
	},
}

var listenPort int

type GuppyStatus struct {
	IsAuthenticated bool   `json:"isAuthenticated"`
	UploadStarted   bool   `json:"uploadStarted"`
	UploadFinished  bool   `json:"uploadFinished"`
	UploadError     string `json:"uploadError"`
	WebLink         string `json:"webLink"`
}

var driveService *drive.Service

func upload(filePath string) {
	// Find an available port
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		slog.Error("Failed to find an available port", "error", err)
		return
	}
	listenPort = listener.Addr().(*net.TCPAddr).Port
	listener.Close()

	// Get the Google OAuth2 client ID from the configuration
	domain := "com.forevanyeung.guppy"
	value := cf.CFPreferencesCopyAppValue("GoogleOauth2ClientId", domain)
	if value == nil {
		slog.Error("Missing GoogleOauth2ClientId configuration")
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
