package main

import (
	"os"

	flags "github.com/jessevdk/go-flags"
	"github.com/monochromegane/goq"
)

var opts goq.Option

func main() {
	args, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(1)
	}

	goq.Query(args[0], args[1])
}
