package main

import (
	"context"
	"encoding/json"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/fzerorubigd/gobgg"
)

func main() {
	var username string
	flag.StringVar(&username, "username", "", "the username")

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGKILL,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGABRT)
	defer cancel()
	flag.Parse()
	bgg := gobgg.NewBGGClient()

	var plays []gobgg.Play
	for i := 0; ; i++ {
		p, err := bgg.Plays(ctx, gobgg.SetUserName(username), gobgg.SetPageNumber(i+1))
		if err != nil {
			panic(err)
		}

		if len(p.Items) == 0 {
			break
		}

		plays = append(plays, p.Items...)
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(plays)
}
