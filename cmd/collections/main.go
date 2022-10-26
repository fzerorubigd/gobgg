package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/fzerorubigd/gobgg"
)

var allType = []gobgg.CollectionType{
	gobgg.CollectionTypeOwn,
	gobgg.CollectionTypeRated,
	gobgg.CollectionTypePlayed,
	gobgg.CollectionTypeComment,
	gobgg.CollectionTypeTrade,
	gobgg.CollectionTypeWant,
	gobgg.CollectionTypeWishList,
	gobgg.CollectionTypePreorder,
	gobgg.CollectionTypeWantToPlay,
	gobgg.CollectionTypeWantToBuy,
	gobgg.CollectionTypePrevOwned,
	gobgg.CollectionTypeHasParts,
	gobgg.CollectionTypeWantParts,
}

func main() {
	var (
		username string
		items    = map[gobgg.CollectionType]*bool{}
	)
	flag.StringVar(&username, "username", "", "the username")
	for _, ct := range allType {
		items[ct] = flag.Bool(string(ct), false, fmt.Sprintf("Include %q items", ct))
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGKILL,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGABRT)
	defer cancel()

	bgg := gobgg.NewBGGClient()

	var opt []gobgg.CollectionType

	for ct, ok := range items {
		if *ok {
			opt = append(opt, ct)
		}
	}
	flag.Parse()
	if username == "" {
		log.Fatal("username is mandatory")
	}

	p, err := bgg.GetCollection(ctx, username, gobgg.SetCollectionTypes(opt...))
	if err != nil {
		log.Fatal(err)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(p)
}
