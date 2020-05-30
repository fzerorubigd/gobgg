package gobgg

import (
	"context"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSearch(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Exact URL match
	httpmock.RegisterResponder("GET", "https://boardgamegeek.com/"+searchPath,
		httpmock.NewStringResponder(200, `<?xml version="1.0" encoding="utf-8"?><items total="7" termsofuse="https://boardgamegeek.com/xmlapi/termsofuse">			<item type="boardgame" id="204583">
			<name type="primary" value="Kingdomino"/>			
							<yearpublished value="2016" />
					</item>
			<item type="boardgame" id="281960">
			<name type="primary" value="Kingdomino Duel"/>			
							<yearpublished value="2019" />
					</item>
			<item type="boardgame" id="260941">
			<name type="primary" value="Kingdomino für 2 Spieler"/>			
							<yearpublished value="2018" />
					</item>
			<item type="boardgame" id="240909">
			<name type="primary" value="Kingdomino: Age of Giants"/>			
							<yearpublished value="2018" />
					</item>
			<item type="boardgame" id="306171">
			<name type="primary" value="Kingdomino: The Court"/>			
							<yearpublished value="2020" />
					</item>
				<item type="boardgameexpansion" id="240909">
			<name type="primary" value="Kingdomino: Age of Giants"/>			
							<yearpublished value="2018" />
					</item>
			<item type="boardgameexpansion" id="306171">
			<name type="primary" value="Kingdomino: The Court"/>			
							<yearpublished value="2020" />
					</item>
						</items>`))

	bgg := NewBGGClient()
	items, err := bgg.Search(context.Background(), "Kingdomino")
	require.NoError(t, err)
	require.Len(t, items, 7)
	assert.Equal(t, "Kingdomino", items[0].Name)
	assert.Equal(t, int64(204583), items[0].ID)
	assert.Equal(t, 2016, items[0].YearPublished)

}
