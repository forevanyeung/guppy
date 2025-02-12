package main

import (
	"github.com/forevanyeung/guppy/cli/analytics"
	"github.com/forevanyeung/guppy/cli/cmd"
)

func main() {
	cmd.Execute()
	analytics.Close()
}
