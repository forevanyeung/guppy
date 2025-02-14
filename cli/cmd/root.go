package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/forevanyeung/guppy/cli/analytics"
	"github.com/forevanyeung/guppy/cli/internal"
	"github.com/spf13/cobra"
)

var verbose bool
var desktop bool

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolVar(&desktop, "desktop", false, "")
	rootCmd.PersistentFlags().MarkHidden("desktop")
	cobra.OnInitialize(initConfig)
}

type warnLevel struct {
	slog.Handler
}

func (warnLevel) Enabled(_ context.Context, level slog.Level) bool {
	return level >= slog.LevelWarn
}

func initConfig() {
	internal.SetVerbose(verbose)
	internal.SetDesktop(desktop)
	analytics.Initialize()

	// unsafe
	if !internal.IsVerbose() {
		h := slog.Default().Handler()
		*slog.Default() = *slog.New(warnLevel{h})
	}
}

var rootCmd = &cobra.Command{
	Use:   "guppy [file]",
	Short: "Guppy is simple tool for opening files in Google Drive",
	Long:  `Guppy can be used as a file handler to associate with file types and open them in Google Drive, or also as
			a command line tool to upload files to Google Drive.`,
	Args:  cobra.MaximumNArgs(1),
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
	Run: func(cmd *cobra.Command, args []string) {
		analytics.TrackEvent("$pageview", map[string]interface{}{
			"$current_url": "root",
		})

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
	if err := rootCmd.Execute(); err != nil {
		slog.Error("Error executing root command", "err", err)
		os.Exit(1)
	}
}
