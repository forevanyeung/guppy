package cmd

import (
	"fmt"

	"github.com/forevanyeung/guppy/analytics"
	"github.com/spf13/cobra"
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
		fmt.Println("Login command executed")
	},
}
