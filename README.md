Boardgamegeek API Client for Golang
===
[![Powered by BoardGameGeek](https://cf.geekdo-images.com/HZy35cmzmmyV9BarSuk6ug__small/img/gbE7sulIurZE_Tx8EQJXnZSKI6w=/fit-in/200x150/filters:strip_icc()/pic7779581.png)](https://boardgamegeek.com)

[![Go Report Card](https://goreportcard.com/badge/github.com/fzerorubigd/gobgg)](https://goreportcard.com/report/github.com/fzerorubigd/gobgg)
[![codecov](https://codecov.io/gh/fzerorubigd/gobgg/branch/master/graph/badge.svg?token=G7VH4EDM0Y)](https://codecov.io/gh/fzerorubigd/gobgg)
[![Go Reference](https://pkg.go.dev/badge/github.com/fzerorubigd/gobgg.svg)](https://pkg.go.dev/github.com/fzerorubigd/gobgg)

Basics
---
Create a new client using the `NewBGGClient`. Providing a bearer token is mandatory.

```go
bgg := gobgg.NewBGGClient(gobgg.SetAuthToken("YOUR_BGG_TOKEN"))
```

You can also set other options, like a custom `http.Client`:
```go
bgg := gobgg.NewBGGClient(
    gobgg.SetAuthToken("YOUR_BGG_TOKEN"),
    gobgg.SetClient(client),
)
```

Collection API
---
Getting the collection in boardgamegeek is a little tricky, so the library will retry with a backoff
to get the collection. You can control the timeout using the context passed to the function 
or it will retry forever. 

```go 
collection, err := bgg.GetCollection(ctx, "fzerorubigd", gobgg.SetCollectionTypes(gobgg.CollectionTypeOwn))
```

Things API
---
You can get the things detail (I just use the board game related API so far) it always 
comes with the statistics (why it should be a flag!)

I decided to follow the same pattern for all the functions even when it is not required. 

```go
things , err := bgg.GetThings(ctx, gobgg.GetThingIDs(1, 2, 3))
```

Plays API
---

Getting a play is like this: 
```go 
// For users
plays, err := bgg.Plays(ctx, gobgg.SetUserName("fzerorubigd"))
// For a game
plays, err := bgg.Plays(ctx, gobgg.SetGameID(1))
```

Posting Plays
---
Posting play is an experimental API that is not using any documented API end point, for this 
you need to call `Login` first. 

Person API
---
For getting the person image you can use the -undocumented- person API `PersonImage`

GeekList
--- 
The library supports getting the lists by their id using the `GeekList` function 

Hotness
--- 
It is possible to get thodays hotness and the the change since yesterday using the `Hotness` function 


Rate Limiting 
---

Boardgamegeek does rate limit on all the api calls. Its a good idea to add a rate limiter to your client. 
it is already supported and I recommend the `go.uber.org/ratelimit`

```go
import (
	"fmt"
	"time"

	"go.uber.org/ratelimit"
	"github.com/fzerorubigd/gobgg"
)

rl := ratelimit.New(10, ratelimit.Per(60*time.Second)) // creates a 10 per minutes rate limiter.
client := gobgg.NewBGGClient(gobgg.SetLimiter(rl))
```
