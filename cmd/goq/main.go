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

	columns, rows := goq.Query(args[0], args[1])

	goban.Render(columns, rows)
}
