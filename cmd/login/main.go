package main

import (
	"context"
	"fmt"
	"github.com/fzerorubigd/gobgg"
)

func main() {
	bgg := gobgg.NewBGGClient()
	bgg.Login(context.Background(), "fzerorubigd", "")

	fmt.Println(bgg.SetRank(context.Background(), 238799, 8.5))
}
