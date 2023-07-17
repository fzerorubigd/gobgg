package gobgg

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func thingResponders(req *http.Request) (*http.Response, error) {
	const head = `<?xml version="1.0" encoding="utf-8"?>
	<items termsofuse="https://boardgamegeek.com/xmlapi/termsofuse">`
	items := func(ids ...int64) string {
		rep := head
		for _, id := range ids {
			year := fmt.Sprintf("%d", 2018+id)
			rep += fmt.Sprintf(`<item type="%[2]s" id="%[1]d">
			<thumbnail>http://thumbnail/%[1]d.png</thumbnail>
			<image>http://image/%[1]d.png</image>
			<name type="primary" value="%[2]s-%[1]d"/>
			<name type="alternate" sortindex="1" value="alter1-%[2]s-%[1]d"/>
			<name type="alternate" sortindex="1" value="alter2-%[2]s-%[1]d"/>
			<description>description for %[2]s-%[1]d</description>
			<yearpublished value="%[3]s" />
			<minplayers value="%[4]d" />
			<maxplayers value="%[5]d" />
			<playingtime value="%[6]d" />
			<minplaytime value="%[7]d" />
			<maxplaytime value="%[8]d" />
			<minage value="%[9]d" />
			<link type="cat1" id="1" value="cat1-v1" />
			<link type="cat1" id="2" value="cat1-v2" />
			<link type="cat1" id="3" value="cat1-v3" />
			<link type="cat2" id="4" value="cat2-v1" />
			<link type="cat2" id="5" value="cat2-v2" />
			<link type="cat3" id="6" value="cat3-v1" />
			</item>`, id, "boardgame", year, id%10+1, id%10+5, id%100+20, id%100-20, id%100+60, id%5+8)
		}
		return rep + "</items>"
	}

	idStr := req.URL.Query().Get("id")
	idArr := strings.Split(idStr, ",")
	ids := make([]int64, 0, len(idArr))
	for i := range idArr {
		id, err := strconv.ParseInt(idArr[i], 10, 0)
		if err == nil {
			ids = append(ids, id)
		}
	}

	return httpmock.NewStringResponse(200, items(ids...)), nil
}

func TestGetThing(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Exact URL match
	httpmock.RegisterResponder(
		"GET",
		"https://boardgamegeek.com/"+thingPath,
		thingResponders,
	)

	bgg := NewBGGClient()
	ids := []int64{1, 2, 3, 4, 5}
	items, err := bgg.GetThings(context.Background(), GetThingIDs(ids...))
	require.NoError(t, err)
	require.Len(t, items, len(ids))
	for i := range ids {
		assert.Equal(t, fmt.Sprintf("boardgame-%d", ids[i]), items[i].Name)

		assert.Equal(t, ids[i], items[i].ID)
		assert.Equal(t, int(2018+ids[i]), items[i].YearPublished)
		assert.Equal(t, BoardGameType, items[i].Type)
		assert.Equal(t, fmt.Sprintf("description for %s-%d", "boardgame", ids[i]), items[i].Description)
		assert.ElementsMatch(t,
			strings.Split(fmt.Sprintf("alter1-%[2]s-%[1]d,alter2-%[2]s-%[1]d", ids[i], "boardgame"), ","),
			items[i].AlternateNames)

		// id%10+1, id%10+5, id%100+20, id%100-20, id%100+60, id%5+8
		assert.Equal(t, int(ids[i]%10+1), items[i].MinPlayers)
		assert.Equal(t, int(ids[i]%10+5), items[i].MaxPlayers)
		assert.Equal(t, fmt.Sprint(ids[i]%100+20), items[i].PlayTime)
		assert.Equal(t, fmt.Sprint(ids[i]%100-20), items[i].MinPlayTime)
		assert.Equal(t, fmt.Sprint(ids[i]%100+60), items[i].MaxPlayTime)
		assert.Equal(t, fmt.Sprint(ids[i]%5+8), items[i].MinAge)

		assert.Len(t, items[i].Links, 3)
		assert.ElementsMatch(t, items[i].Links["cat1"], []Link{
			{ID: 1, Name: "cat1-v1"},
			{ID: 2, Name: "cat1-v2"},
			{ID: 3, Name: "cat1-v3"},
		})
		assert.ElementsMatch(t, items[i].Links["cat2"], []Link{
			{ID: 4, Name: "cat2-v1"},
			{ID: 5, Name: "cat2-v2"},
		})
		assert.ElementsMatch(t, items[i].Links["cat3"], []Link{
			{ID: 6, Name: "cat3-v1"},
		})
	}
}
