package analytics

import (
	"fmt"
	"log/slog"

	"github.com/denisbrodbeck/machineid"
	"github.com/posthog/posthog-go"
)

var PosthogApiKey = ""
var PosthogEndpoint = ""

var client posthog.Client
var id string

func init() {
	client, _ = posthog.NewWithConfig(
		PosthogApiKey,
		posthog.Config{
			Endpoint: PosthogEndpoint,
			Verbose: false,
		},
	)

	id, _ = machineid.ProtectedID("guppy")

	slog.Info("Analytics initialized", "Machine Id", id)
}

func TrackEvent(event string, properties map[string]interface{}) {
	// TODO: Add a configuration to disable analytics
	// disableAnalytics := cf.CFPreferencesCopyAppValue("DisableAnalytics", "com.forevanyeung.guppy")
	// if disableAnalytics != nil && disableAnalytics.(bool) {
	// 	return
	// }

	if client == nil {
		return
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
