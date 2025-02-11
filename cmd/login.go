package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(loginCmd)
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to Google Drive",
	Run: func(cmd *cobra.Command, args []string) {
		// Implement login logic here
		fmt.Println("Login command executed")
	},
}
