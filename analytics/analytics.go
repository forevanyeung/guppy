package analytics

import (
	"fmt"

	"github.com/denisbrodbeck/machineid" // Import the cf package
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
		},
	)

	id, _ = machineid.ProtectedID("guppy")

	fmt.Println("Machine ID:", id)
}

func TrackEvent(event string, properties map[string]interface{}) {
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

	fmt.Println("Event tracked:", event)
}

func Close() {
	if client == nil {
		return
	}

	client.Close()
}
