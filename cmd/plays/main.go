package main

import (
	"encoding/json"
	"flag"
	"os"

	"github.com/fzerorubigd/clictx"

	"github.com/fzerorubigd/gobgg"
)

func main() {
	var (
		pass string
	)

	flag.StringVar(&pass, "password", "", "the pass")

	ctx := clictx.DefaultContext()
	bgg := gobgg.NewBGGClient()

	var plays []gobgg.Play
	for i := 0; ; i++ {
		p, err := bgg.Plays(ctx, gobgg.SetUserName("fzerorubigd"), gobgg.SetPageNumber(i+1))
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
