package cmd

import (
	"fmt"
	"log/slog"
	"mime"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/forevanyeung/guppy/cli/analytics"
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

	f := &drive.File{
		Name:     filepath.Base(filePath),
		MimeType: mapMimeTypeToGoogleMimeType(mimeType),
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

	// Calculate file size after successful upload
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		slog.Error("Error getting file info:", "err", err)
		// Don't return here as the upload was successful
	}
	fileSizeKB := float64(fileInfo.Size()) / 1024.0

	analytics.TrackEvent("upload_done", map[string]interface{}{
		"auth": "cached|new",
		"duration": 0,
		"file_size_kb": fileSizeKB,
		"file_type": mimeType,
	})
}

func getMimeType(filePath string) string {
	// Get the file extension
	ext := filepath.Ext(filePath)

	// Use the mime package to detect the MIME type
	mimeType := mime.TypeByExtension(ext)
	return mimeType
}

func mapMimeTypeToGoogleMimeType(mimeType string) string {
	// Map the MIME type to the Google MIME type
	switch mimeType {
	case "text/csv":
		return "application/vnd.google-apps.spreadsheet"
	case "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":
		return "application/vnd.google-apps.spreadsheet"
	case "text/plain":
		return "application/vnd.google-apps.document"
	case "application/pdf":
		return "application/vnd.google-apps.document"
	case "application/vnd.openxmlformats-officedocument.wordprocessingml.document":
		return "application/vnd.google-apps.document"
	case "application/vnd.openxmlformats-officedocument.presentationml.presentation":
		return "application/vnd.google-apps.presentation"
	default:
		return "application/vnd.google-apps.unknown"
	}
}
