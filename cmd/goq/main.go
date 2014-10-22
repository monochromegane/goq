package main

import (
	"os"

	flags "github.com/jessevdk/go-flags"
	"github.com/monochromegane/goban"
	"github.com/monochromegane/goq"
)

var opts goq.Option

func main() {
	args, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(1)
	}

	if opts.List {
		goq.List(args[0])
		os.Exit(0)
	}

	columns, rows := goq.Query(args[0], args[1], args[2:]...)

	goban.Render(columns, rows)
}
