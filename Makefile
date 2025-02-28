#!make

ifneq (,$(wildcard .env))
    include .env
endif

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

.PHONY: clean

clean:
	rm -rf ./cli/dist

.PHONY: build

build: clean
	go build -ldflags " \
	-X github.com/forevanyeung/guppy/cli/analytics.PosthogEndpoint=$(POSTHOG_ENDPOINT) \
	-X github.com/forevanyeung/guppy/cli/analytics.PosthogApiKey=$(POSTHOG_API_KEY) \
	$(if $(VERSION),-X github.com/forevanyeung/guppy/cli/internal.Version=$(VERSION))" \
	-o ./cli/dist/guppy ./cli
