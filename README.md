# Recommender
*Recommendation engine written in Go.*

**This repository is intended to be an exercise, not a package for any production code.**

## Why?

I've built this project to learn about, among other things, recommendation basics, concurrency patterns, testing techniques, and key/value persistent storage. Should you browse the code, you'll find a few examples of each.

## Examples

Although the best examples are documented in the test package, here are some brief examples of how the system works:

```go
// Create a user
niko := recommender.NewUser("Niko Kovacevic")

// Create items
flagstaff := recommender.NewItem("Flagstaff, Arizona")
losAngeles := recommender.NewItem("Los Angeles, California")

// Create a recommender
r := recommender.NewRecommender()

// Rate items
r.Like(niko, flagstaff)
r.Dislike(niko, losAngeles)

// Updating happens automatically upon rating an item.
// Get some suggestions!
suggestions, err = r.GetSuggestions(niko)
if err != nil {
  return err
}

```

If other people rated your items, as well as some other items, have a look at the suggestions, scored from -1 to 1.

```bash
{Portland, Oregon 0.38888893}
{San Francisco, California -0.1111111}
{Sacramento, California -0.11111113}
{Princeton, New Jersey -1}
{Austin, Texas 0.5}
{Philadelphia, Pennsylvania -1}
{Ashland, Oregon 0.44444445}
{Tacoma, Washington 0.3333333}
{Houston, Texas -0.5}
{Tucson, Arizona 1}
{New York, New York 0}
{Santa Fe, New Mexico 0.8333334}
{Portland, Maine -0.3333333}
```

## Next

I haven't determined whether or not to extend the project by building a front-end. Were I to go that direction, I'd likely build a React application with OAuth (perhaps leveraging Auth0) to let users sign in and rate cities or movies via a Go API, which would leverage this package.

## Acknowlegements

- https://www.toptal.com/algorithms/predicting-likes-inside-a-simple-recommendation-engine

## License

The MIT License (MIT)  
Copyright (c) 2016 Niko Kovacevic
