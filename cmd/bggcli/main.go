package main

import (
	"context"
	"flag"
	"log"
	"os/signal"
	"syscall"

	"github.com/fzerorubigd/gobgg"
)

func main() {
	ctx, cnl := signal.NotifyContext(
		context.Background(),
		syscall.SIGKILL,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGABRT,
	)
	defer cnl()

	// TODO: Add option to customize the http client
	flag.Usage = usage
	flag.Parse()

	bgg := gobgg.NewBGGClient()

	if err := dispatch(ctx, bgg, flag.Args()...); err != nil {
		log.Fatal(err.Error())
	}
}
