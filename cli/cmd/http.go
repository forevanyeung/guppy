package cmd

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
)

type OAuthResponse struct {
	State       string `json:"state"`
	AccessToken string `json:"accessToken"`
	TokenType   string `json:"tokenType"`
	ExpiresIn   int16  `json:"expiresIn"`
	Scope       string `json:"scope"`
}

type OAuthResponseError struct {
	Error string `json:"error"`
}

func httpServer(newAuthChan chan string, statusChan chan GuppyStatus, serverDoneChan chan bool) {
	// Get the directory of the current file
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)

	// Construct the path relative to the current file
	staticDir := filepath.Join(dir, "../static")
	
	mux := http.NewServeMux()
	mux.Handle("GET /", http.FileServer(http.Dir(staticDir)))
	mux.HandleFunc("POST /auth", authHttp(newAuthChan))
	mux.HandleFunc("GET /status", statusHttp(statusChan, serverDoneChan))

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", listenPort),
		Handler: mux,
	}

	// Start the server in a goroutine
	go func() {
		slog.Info(fmt.Sprintf("Listening on port %d", listenPort))

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server error", "err", err)
		}
	}()

	// Handle graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// Block until a signal is received
	<-c

	slog.Info("Shutting down server...")
	server.Shutdown(nil) // Gracefully shut down the server
	slog.Info("Server stopped")
}

func authHttp(newAuthChan chan string) func (w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Body == nil {
			http.Error(w, "Please send a request body", http.StatusBadRequest)
			return
		}

		// decode the json in body
		var res OAuthResponse
		err := json.NewDecoder(r.Body).Decode(&res)
		if err != nil {
			http.Error(w, "Unable to decode request", http.StatusBadRequest)
			slog.Error("Error decoding request", "err", err)
			return
		}

		// TODO validate the state

		// Send the new auth token through the channel
		newAuthChan <- res.AccessToken

		return
	}
}

func statusHttp(statusChan chan GuppyStatus, serverDoneChan chan bool) func (w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO implement long polling

		status := <-statusChan

		jsonResponse, err := json.Marshal(status)
		if err != nil {
			http.Error(w, "Unable to encode response", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)

		// only close the server once it sent an upload finished status
		if status.UploadFinished {
			serverDoneChan <- true
		}
	}
}
