package main

import (
	"flag"
	"log"
	"syscall"

	"github.com/fzerorubigd/clictx"
	"github.com/fzerorubigd/gobgg"
)

func main() {
	ctx := clictx.Context(
		syscall.SIGKILL,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGABRT,
	)

	// TODO: Add option to customize the http client
	flag.Usage = usage
	flag.Parse()

	bgg := gobgg.NewBGGClient()

	if err := dispatch(ctx, bgg, flag.Args()...); err != nil {
		log.Fatal(err.Error())
	}
}
