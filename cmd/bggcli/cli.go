package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"sync"

	"github.com/fzerorubigd/gobgg"
)

var (
	commands []command
	lock     sync.RWMutex
)

type command struct {
	Name        string
	Description string

	Run func(context.Context, *gobgg.BGG, ...string) error
}

func usage() {
	// usage is also called in dispatch, but multiple read lock is fine
	lock.RLock()
	defer lock.RUnlock()

	fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()

	fmt.Fprintf(flag.CommandLine.Output(), "Sub commands:\n")

	for i := range commands {
		fmt.Fprintf(flag.CommandLine.Output(), "  %s: %s\n", commands[i].Name, commands[i].Description)
	}
}

func dispatch(ctx context.Context, bgg *gobgg.BGG, args ...string) error {
	lock.RLock()
	defer lock.RUnlock()

	if len(args) < 1 {
		return fmt.Errorf("atleast one arg is required")
	}

	sub := args[0]
	for i := range commands {
		if sub == commands[i].Name {
			return commands[i].Run(ctx, bgg, args...)
		}
	}

	usage()
	return fmt.Errorf("invalid command")
}

func addCommand(name, description string, run func(context.Context, *gobgg.BGG, ...string) error) {
	lock.Lock()
	defer lock.Unlock()

	commands = append(commands, command{
		Name:        name,
		Description: description,
		Run:         run,
	})
}
