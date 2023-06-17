package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/fzerorubigd/gobgg"
)

func thing(ctx context.Context, bgg *gobgg.BGG, args ...string) error {
	cmd := flag.NewFlagSet(args[0], flag.ExitOnError)
	var (
		links       bool
		names       bool
		description bool
	)

	cmd.BoolVar(&links, "links", false, "Show links")
	cmd.BoolVar(&names, "names", false, "Show alternate names")
	cmd.BoolVar(&description, "description", true, "Show description")
	if err := cmd.Parse(args[1:]); err != nil {
		return err
	}

	idStr := strings.Join(cmd.Args(), " ")
	id, err := strconv.ParseInt(idStr, 10, 0)
	if err != nil {
		return fmt.Errorf("the argument should be an integer, but is %q", idStr)
	}

	opts := []gobgg.GetOptionSetter{
		gobgg.GetThingIDs(id),
	}
	result, err := bgg.GetThings(ctx, opts...)
	if err != nil {
		return err
	}

	for _, item := range result {
		year := "Not specified"
		if item.YearPublished > 0 {
			year = fmt.Sprint(item.YearPublished)
		}
		fmt.Fprintf(os.Stdout, "%d\t%s\n\n", item.ID, item.Name)
		fmt.Fprintf(os.Stdout, "%s, Published in %s\n", item.Type, year)
		fmt.Fprintf(os.Stdout, "Play time: %s (Min: %s, Max:%s)\n", item.PlayTime, item.MinPlayTime, item.MaxPlayTime)
		if item.MinPlayers != item.MaxPlayers {
			fmt.Fprintf(os.Stdout, "Players count: %d-%d\n", item.MinPlayers, item.MaxPlayers)
		} else {
			fmt.Fprintf(os.Stdout, "Players count: %d\n", item.MinPlayers)
		}
		fmt.Fprintln(os.Stdout, "Suggested Player count (community votes): ")
		for i := range item.SuggestedPlayerCount {
			rec, num, per := item.SuggestedPlayerCount[i].Suggestion()
			fmt.Fprintf(os.Stdout, "%s => %s, %d votes, %0.2f%%\n",
				item.SuggestedPlayerCount[i].NumPlayers,
				rec, num, per)
		}
		if names && len(item.AlternateNames) > 0 {
			fmt.Fprintf(os.Stdout, "Alternate names: %s\n\n", strings.Join(item.AlternateNames, ", "))
		}
		if description {
			fmt.Fprintln(os.Stdout, item.Description)
		}
		if links {
			keys := make([]string, 0, len(item.Links))
			for i := range item.Links {
				keys = append(keys, i)
			}
			sort.Strings(keys)
			for _, key := range keys {
				fmt.Fprintf(os.Stdout, "%s: \n", key)
				for _, lnk := range item.Links[key] {
					fmt.Printf("\t%d: %s\n", lnk.ID, lnk.Name)
				}
			}
		}
	}

	return nil
}

func init() {
	addCommand("thing", "Get a thing from bgg", thing)
}
