package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"
	"strings"

	"github.com/forevanyeung/guppy/cli/analytics"
	"github.com/spf13/cobra"
	"github.com/zalando/go-keyring"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

func init() {
	rootCmd.AddCommand(loginCmd)
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to Google Drive",
	Run: func(cmd *cobra.Command, args []string) {
		analytics.TrackEvent("$pageview", map[string]interface{}{
			"$current_url": "login",
		})

		// Implement login logic here
		slog.Info("Login command executed")
	},
}

type Token struct {
	*oauth2.Token
}

type OAuthConfig struct {
	*oauth2.Config
}

func auth(newAuthChan chan string, clientId string) {
	var token Token
	var pass bool

	// Check keyring if an auth token exists
	auth, err := keyring.Get("com.forevanyeung.guppy.auth", "AccessToken")
	if err == nil {
		slog.Info("Retrieved auth token from keyring")

		token = Token{
			&oauth2.Token{AccessToken: auth},
		}

		if token.IsValid() {
			pass = true	
			openBrowser(fmt.Sprintf("http://localhost:%d/interstitial.html", listenPort))
		} else {
			slog.Info("Auth token from keyring is no longer valid")
		}
	}

	// If no token is found, or token is invalid, get a new token
	if !pass {
		slog.Info("Getting a new auth token, opening browser")

		config := OAuthConfig{
			Config: &oauth2.Config{
				ClientID:    clientId,
				RedirectURL: fmt.Sprintf("http://localhost:%d/interstitial.html", listenPort),
				Scopes:      []string{"https://www.googleapis.com/auth/drive"},
				Endpoint:    google.Endpoint,
			},
		}

		config.Login()

		// wait for response from auth server
		newAuth := <-newAuthChan
		slog.Info("Received a new auth token")

		// Save the token in keyring
		err = keyring.Set("com.forevanyeung.guppy.auth", "AccessToken", newAuth)
		if err != nil {
			slog.Error("Error saving token to keyring", "err", err)
			return
		}

		// set token to use
		token = Token {
			&oauth2.Token{AccessToken: newAuth},
		}
	}

	// Create a new Drive service with the auth token
	ctx := context.Background()
	driveService, err = drive.NewService(ctx, option.WithTokenSource(oauth2.StaticTokenSource(token.Token)))
	if err != nil {
		slog.Error("Error creating Drive service", "err", err)
		return
	}
}

func (c *OAuthConfig) Login() {
	url := c.AuthCodeURL(
		// TODO random generate state
		"state",
		oauth2.SetAuthURLParam("response_type", "token"),
	)
	
	openBrowser(url)
}

func (t *Token) Logout() {
	keyring.Delete("com.forevanyeung.guppy.auth", "AccessToken")
	slog.Info("Logged out and token deleted from keyring")
}

func (t *Token) IsValid() bool {
	uri := fmt.Sprintf("https://www.googleapis.com/oauth2/v3/tokeninfo?access_token=%s", t.Token.AccessToken)
	res, err := http.Get(uri)
	if err != nil {
		slog.Error("Error checking token validity", "err", err)
		return false
	}

	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return false
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		slog.Error("Error decoding JSON response", "err", err)
		return false
	}

	// TODO: Check if token has the required scopes
	// tokenScopes := result["scope"].(string)
	// for _, s := range scope {
	// 	if !contains(tokenScopes, s) {
	//		slog.Warn("Token does not have the required scope")
	// 		return false
	// 	}
	// }

	return true
}

func contains(scopes string, scope string) bool {
	for _, s := range strings.Split(scopes, " ") {
		if s == scope {
			return true
		}
	}
	return false
}

// This generates a 16-character state string (128 bits), which is generally sufficient for local OAuth flows.
// If you want it even simpler and faster, you can use Go’s math/rand (less secure but faster).
func generateState() string {
	return fmt.Sprintf("%016x", rand.Int63())
}
