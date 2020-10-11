package main

import (
	"flag"
	"log"
	"syscall"

	"github.com/fzerorubigd/clictx"

	"github.com/fzerorubigd/gobgg"
)

func main() {
	var (
		pass, user string
	)

	flag.StringVar(&pass, "password", "", "the pass")
	flag.StringVar(&user, "username", "", "the username")
	flag.Parse()

	ctx := clictx.Context(syscall.SIGTERM)
	bgg := gobgg.NewBGGClient()

	err := bgg.Login(ctx, user, pass)
	if err != nil {
		panic(err)
	}

	for i := 0; ; i++ {
		p, err := bgg.Plays(ctx, gobgg.SetUserName("fzerorubigd"), gobgg.SetPageNumber(i+1))
		if err != nil {
			panic(err)
		}

		if len(p.Items) == 0 {
			break
		}

		for _, pl := range p.Items {
			if err := bgg.PostPlay(ctx, pl); err != nil {
				log.Fatal(err)
			}
		}
	}
}
