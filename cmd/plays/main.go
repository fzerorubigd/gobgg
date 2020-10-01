package main

import (
	"encoding/json"
	"fmt"
	"github.com/fzerorubigd/clictx"
	"github.com/fzerorubigd/gobgg"
	"syscall"
)

func main() {
	ctx := clictx.Context(syscall.SIGTERM)
	bgg := gobgg.NewBGGClient()
	var all *gobgg.PlaysFixed
	for i := 0; ; i++ {
		p, err := bgg.Plays(ctx, gobgg.SetUserName("fzerorubigd"), gobgg.SetPageNumber(i+1))
		if err != nil {
			panic(err)
		}

		if len(p.Items) == 0 {
			break
		}

		if all == nil {
			all = p
		}
		all.Items = append(all.Items, p.Items...)
	}

	b, err := json.MarshalIndent(all, "", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Print(string(b))
}
