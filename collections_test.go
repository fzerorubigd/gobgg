package gobgg_test

import (
	"context"
	"testing"

	"github.com/fzerorubigd/gobgg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetCollection(t *testing.T) {
	ctx := context.Background()

	bgg := gobgg.NewBGGClient()
	col, err := bgg.GetCollection(ctx, "gobgg")
	require.NoError(t, err)
	// I create the user and add some games to it
	games := map[int64][]string{
		161936: {"previouslyowned"},      // Pandemic Legacy: Season 1
		174430: {"wishlist", "musthave"}, // Gloomhaven
		204583: {"owned"},
		224517: {"owned"},
		23383:  {"owned", "played"},
		233078: {"wishlist", "thinkingaboutit"}, // TI4
		342942: {"wishlist", "liketohave"},      // Ark Nova
	}

	assert.Equal(t, len(games), len(col))
	for i := range col {
		require.ElementsMatch(t, games[col[i].ID], col[i].CollectionStatus)
	}
}
