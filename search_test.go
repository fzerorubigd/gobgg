package gobgg

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func serachResponders(req *http.Request) (*http.Response, error) {
	const head = `<?xml version="1.0" encoding="utf-8"?>
	<items total="%d" termsofuse="https://boardgamegeek.com/xmlapi/termsofuse">`

	items := func(name string, count int, types string) string {
		if types == "" {
			types = "rpgitem,videogame,boardgame,boardgameaccessory,boardgameexpansion"
		}

		tpArray := strings.Split(types, ",")

		rep := fmt.Sprintf(head, count)
		for id := 0; id < count; id++ {
			tp := tpArray[rand.Intn(len(tpArray))]
			year := "2018"
			if id%3 == 2 {
				year = ""
			}
			rep += fmt.Sprintf(`<item type="%s" id="%d">
			<name type="primary" value="%s-%d"/>
			<yearpublished value="%s" />
			</item>`, tp, id+1000, name, id, year)
		}
		return rep + "</items>"
	}

	query := req.URL.Query()
	search := query.Get("query")
	if search == "empty" {
		return httpmock.NewStringResponse(200, items(search, 0, "")), nil
	}

	exact := query.Get("exact")
	count := 10
	if exact == "1" {
		count = 1
	}

	return httpmock.NewStringResponse(200, items(search, count, query.Get("type"))), nil
}

func TestSearch(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Exact URL match
	httpmock.RegisterResponder(
		"GET",
		"https://boardgamegeek.com/"+searchPath,
		serachResponders,
	)

	bgg := NewBGGClient()
	items, err := bgg.Search(context.Background(), "Kingdomino", SearchExact(), SearchTypes(BoardGameType))
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "Kingdomino-0", items[0].Name)
	assert.Equal(t, int64(1000), items[0].ID)
	assert.Equal(t, 2018, items[0].YearPublished)
	assert.Equal(t, BoardGameType, items[0].Type)
}
