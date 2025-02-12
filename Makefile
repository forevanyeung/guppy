#!make
include .env

.PHONY: run

run:
	go run -ldflags " \
	-X github.com/forevanyeung/guppy/cli/analytics.PosthogEndpoint=$(POSTHOG_ENDPOINT) \
	-X github.com/forevanyeung/guppy/cli/analytics.PosthogApiKey=$(POSTHOG_API_KEY)" \
	./cli $(FILE) $(if $(VERBOSE),-v)

.PHONY: run-desktop

run-desktop:
	go run -ldflags " \
	-X github.com/forevanyeung/guppy/cli/analytics.PosthogEndpoint=$(POSTHOG_ENDPOINT) \
	-X github.com/forevanyeung/guppy/cli/analytics.PosthogApiKey=$(POSTHOG_API_KEY)" \
	./cli $(FILE) --desktop $(if $(VERBOSE),-v)
