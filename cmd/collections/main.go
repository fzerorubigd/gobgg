package main

import (
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
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

func getTheUrl(ctx context.Context, id int64) string {
	result := fmt.Sprintf("https://boardgamegeek.com/boardgame/%d", id)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, result, nil)
	if err != nil {
		return result
	}
	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		return result
	}
	r := resp.Header.Get("Location")
	if r != "" {
		result = "https://boardgamegeek.com" + r
	}

	return result
}

func main() {
	var (
		username string
		items    = map[gobgg.CollectionType]*bool{}
	)
	flag.StringVar(&username, "username", "fzerorubigd", "the username")
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

	wcsv := csv.NewWriter(os.Stdout)
	defer wcsv.Flush()
	for i := range p {
		rec := []string{
			p[i].Name,
			fmt.Sprint(p[i].ID),
			getTheUrl(ctx, p[i].ID),
			strings.Join(p[i].CollectionStatus, ","),
		}

		_ = wcsv.Write(rec)
	}
}
