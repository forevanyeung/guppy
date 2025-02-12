#!make
include .env

.PHONY: run

run:
	go run -ldflags "-X github.com/forevanyeung/guppy/analytics.PosthogEndpoint=$(POSTHOG_ENDPOINT) -X github.com/forevanyeung/guppy/analytics.PosthogApiKey=$(POSTHOG_API_KEY)" . $(FILE)
