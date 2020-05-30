package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/fzerorubigd/gobgg"
)

func search(ctx context.Context, bgg *gobgg.BGG, args ...string) error {
	cmd := flag.NewFlagSet(args[0], flag.ExitOnError)
	var (
		exact bool
	)
	cmd.BoolVar(&exact, "exact", false, "exact search on bgg")
	if err := cmd.Parse(args[1:]); err != nil {
		return err
	}

	opts := []gobgg.SearchOptionSetter{}
	if exact {
		opts = append(opts, gobgg.SearchExact())
	}
	result, err := bgg.Search(ctx, strings.Join(cmd.Args(), " "), opts...)
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 8, 4, 4, ' ', 0)
	defer w.Flush()
	fmt.Fprintln(w, "ID\tName\tType\tYear Published")
	for _, item := range result {
		year := "Not specified"
		if item.YearPublished > 0 {
			year = fmt.Sprint(item.YearPublished)
		}
		fmt.Fprintf(w, "%d\t%s\t%s\t%s\n", item.ID, item.Name, item.Type, year)
	}

	return nil
}

func init() {
	addCommand("search", "Search for a game in boardgae geek", search)
}
