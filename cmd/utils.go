package cmd

import (
	"fmt"
	"log/slog"
	"math/rand"
	"mime"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/forevanyeung/guppy/analytics"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
)

func openBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}

	if err != nil {
		fmt.Println("Please open the URL in your browser:", url)
		slog.Error("Failed to open browser:", "err", err)
	}
}

// This generates a 16-character state string (128 bits), which is generally sufficient for local OAuth flows.
// If you want it even simpler and faster, you can use Go’s math/rand (less secure but faster).
func generateState() string {
	return fmt.Sprintf("%016x", rand.Int63())
}

func uploadFile(filePath string, uploadStatus chan GuppyStatus) {
	file, err := os.Open(filePath)
	if err != nil {
		slog.Error("Error opening file:", "err", err)
		uploadStatus <- GuppyStatus{UploadError: "Error opening file"}
		return
	}
	defer file.Close()

	// get the MIME type of the file
	mimeType := getMimeType(filePath)

	slog.Info(fmt.Sprintf("MIME type: %s", mimeType))

	// TODO: pick mimetype based on file extension
	f := &drive.File{
		Name:     filepath.Base(filePath),
		MimeType: "application/vnd.google-apps.spreadsheet",
	}

	createdFile, err := driveService.Files.Create(f).Media(file, googleapi.ContentType(mimeType)).Fields("webViewLink").Do()
	if err != nil {
		slog.Error("Error creating file:", "err", err)
		uploadStatus <- GuppyStatus{UploadError: "Error uploading file to Google Drive"}
		return
	}

	uploadStatus <- GuppyStatus{
		UploadFinished: true, 
		WebLink: createdFile.WebViewLink,
	}

	// TODO: track actual data for analytics
	analytics.TrackEvent("upload_done", map[string]interface{}{
		"guppy_version": "0.0.0",
		"guppy_platform": "cli",
		"os_platform": "darwin",
		"duration": 0,
		"file_size_kb": 100,
		"file_type": "csv",
	})
}

func getMimeType(filePath string) string {
	// Get the file extension
	ext := filepath.Ext(filePath)

	// Use the mime package to detect the MIME type
	mimeType := mime.TypeByExtension(ext)
	return mimeType
}
