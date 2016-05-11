package main

import "github.com/nikovacevic/recommender"

func main() {
	likeRater := recommender.NewRater("like")
	defer likeRater.CloseDB()

	dislikeRater := recommender.NewRater("dislike")
	defer dislikeRater.CloseDB()
}
