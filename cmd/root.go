package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "guppy [file]",
	Short: "Guppy is simple tool for opening files in Google Drive",
	Long:  `Guppy can be used as a file handler to associate with file types and open them in Google Drive, or also as
			a command line tool to upload files to Google Drive.`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Check if a file path is provided as a command line argument
		if len(args) == 1 {
			filePath := args[0]
			upload(filePath)
		} else {
			fmt.Println("Please provide a file path as a command line argument")
		}
	},
}

func Execute() {
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(uploadCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
