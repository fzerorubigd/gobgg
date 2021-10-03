Boardgamegeek API Client for Golang
===

Basics
---
Create a new client using the `NewBGGClient` also you can use the `Login` function 
to set the cookies, some API (like posting plays) needs the login cookies. 

You can set the `http.Client`, or `http.Cookie` or using different domain (do you need to?)

```go
bgg := gobgg.NewBGGClient(gobgg.SetClient(client))
// Setting the cookies (not a requirement if you want to only use public API)
if err := bgg.Login("user", "password"); err != nil {
	log.Fatalf("Login failed with error: %q", err)	
}
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