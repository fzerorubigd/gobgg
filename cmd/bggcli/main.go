package main

import (
	"context"
	"log"

	"github.com/davecgh/go-spew/spew"
	"github.com/fzerorubigd/gobgg"
)

func main() {
	bgg := gobgg.NewBGGClient()

	items, err := bgg.Search(context.Background(), "Hokm")
	if err != nil {
		log.Fatal(err)
	}

	spew.Dump(items)
}
