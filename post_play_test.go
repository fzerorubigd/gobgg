package gobgg_test

import (
	"context"
	"testing"
	"time"

	"github.com/fzerorubigd/gobgg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostPlay(t *testing.T) {
	//pass := os.Getenv("BGG_PASSWORD")
	pass := "a+xE?t-I>jlOc8]!Ins'9OxtX"
	if pass == "" {
		t.Skip()
		return
	}
	ctx := context.Background()
	bgg := gobgg.NewBGGClient()
	err := bgg.Login(ctx, "gobgg", pass)
	require.NoError(t, err)
	assert.Equal(t, "gobgg", bgg.GetActiveUsername())

	err = bgg.PostPlay(ctx, &gobgg.Play{
		Date:     time.Now(),
		Quantity: 1,
		Length:   time.Minute * 20,
		Location: "Testing",
		Comment:  "Testing",
		Item: gobgg.Item{
			ID:   23383, // Hokm ,
			Type: "thing",
		},
		Players: []gobgg.Player{
			{
				UserName: "fzerorubigd",
				Name:     "Forud",
				Win:      true,
			}, {
				Name:     "GoBGG",
				UserName: "gobgg",
				UserID:   "3597059",
				Win:      false,
			},
		},
	})

	require.NoError(t, err)
}
