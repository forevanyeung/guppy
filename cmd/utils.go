package cmd

import (
	"fmt"
	"math/rand"
	"mime"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

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
		fmt.Printf("Failed to open browser: %v\n", err)
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
		fmt.Println("Error opening file:", err)
		uploadStatus <- GuppyStatus{UploadError: "Error opening file"}
		return
	}
	defer file.Close()

	// get the MIME type of the file
	mimeType := getMimeType(filePath)

	fmt.Println("MIME type:", mimeType)

	f := &drive.File{
		Name:     filepath.Base(filePath),
		MimeType: "application/vnd.google-apps.spreadsheet",
	}

	createdFile, err := driveService.Files.Create(f).Media(file, googleapi.ContentType(mimeType)).Fields("webViewLink").Do()
	if err != nil {
		fmt.Println("Error creating file:", err)
		uploadStatus <- GuppyStatus{UploadError: "Error uploading file to Google Drive"}
		return
	}

	uploadStatus <- GuppyStatus{
		UploadFinished: true, 
		WebLink: createdFile.WebViewLink,
	}
}

func getMimeType(filePath string) string {
	// Get the file extension
	ext := filepath.Ext(filePath)

	// Use the mime package to detect the MIME type
	mimeType := mime.TypeByExtension(ext)
	return mimeType
}
