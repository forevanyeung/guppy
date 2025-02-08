package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
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
	mux := http.NewServeMux()
	mux.Handle("GET /", http.FileServer(http.Dir("./static")))
	mux.HandleFunc("POST /auth", authHttp(newAuthChan))
	mux.HandleFunc("GET /status", statusHttp(statusChan, serverDoneChan))

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", listenPort),
		Handler: mux,
	}

	// Start the server in a goroutine
	go func() {
		fmt.Println("Listening on port", listenPort)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Server error: %v\n", err)
		}
	}()

	// Handle graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// Block until a signal is received
	<-c

	fmt.Println("Shutting down server...")
	server.Shutdown(nil) // Gracefully shut down the server
	fmt.Println("Server stopped")
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
			fmt.Println(err)
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
