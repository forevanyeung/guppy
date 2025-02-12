package main

import (
	"github.com/forevanyeung/guppy/analytics"
	"github.com/forevanyeung/guppy/cmd"
)

func main() {
	cmd.Execute()
	analytics.Close()
}
