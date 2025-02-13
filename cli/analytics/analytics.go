package analytics

import (
	"fmt"
	"log/slog"
	"runtime"

	"github.com/denisbrodbeck/machineid"
	"github.com/forevanyeung/guppy/cli/cf"
	"github.com/forevanyeung/guppy/cli/internal"
	"github.com/posthog/posthog-go"
)

var PosthogApiKey = ""
var PosthogEndpoint = ""

var client posthog.Client
var id string
var meta map[string]interface{}

func Initialize() {
	disableAnalytics := cf.CFPreferencesCopyAppValue("DisableAnalytics", "com.forevanyeung.guppy")
	if disableAnalytics != nil && disableAnalytics.(bool) {
		return
	}

	client, _ = posthog.NewWithConfig(
		PosthogApiKey,
		posthog.Config{
			Endpoint: PosthogEndpoint,
			Verbose: internal.IsVerbose(),
		},
	)

	id, _ = machineid.ProtectedID("guppy")

	meta = map[string]interface{}{
		"guppy_version": internal.Version,
		"guppy_platform": func() string {
			if internal.IsDesktop() {
				return "desktop"
			}
			return "cli"
		}(),
		"os_platform": runtime.GOOS,
	}

	slog.Info("Analytics initialized", "Machine Id", id)
}

func TrackEvent(event string, properties map[string]interface{}) {
	if client == nil {
		return
	}

	// Merge meta properties to be included with every event into the properties map
	for key, value := range meta {
		if _, exists := properties[key]; !exists {
			properties[key] = value
		}
	}

	c := posthog.Capture{
		DistinctId: id,
		Event:      event,
		Properties: properties,
	}

	client.Enqueue(c)

	slog.Info(fmt.Sprintf("Event tracked: %s", event))
}

func Close() {
	if client == nil {
		return
	}

	client.Close()
}
