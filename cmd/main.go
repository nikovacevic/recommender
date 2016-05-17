package main

import "github.com/nikovacevic/recommender"

func main() {
	rater := recommender.NewRater()

	niko := recommender.NewUser("Niko Kovacevic")
	aubreigh := recommender.NewUser("Aubreigh Brunschwig")
	nick := recommender.NewUser("Nick Evers")
	katie := recommender.NewUser("Katie Yoder")
	johnny := recommender.NewUser("Johnny Bernard")
	amanda := recommender.NewUser("Amanda Hunt")

	denver := recommender.NewItem("Denver")
	phoenix := recommender.NewItem("Phoenix")
	portland := recommender.NewItem("Portland")
	seattle := recommender.NewItem("Seattle")
	pittsburgh := recommender.NewItem("Pittsburgh")
	houston := recommender.NewItem("Houston")
	austin := recommender.NewItem("Austin")

	rater.AddLike(niko, phoenix)
	rater.AddLike(aubreigh, phoenix)
	rater.AddLike(johnny, phoenix)
	rater.AddLike(amanda, phoenix)

	rater.AddLike(nick, houston)
	rater.AddLike(katie, houston)
	rater.AddLike(nick, austin)
	rater.AddLike(katie, austin)

	rater.AddLike(nick, pittsburgh)
	rater.AddLike(niko, pittsburgh)

	rater.AddLike(niko, portland)
	rater.AddLike(aubreigh, portland)
	rater.AddLike(nick, portland)
	rater.AddLike(katie, portland)

	rater.AddLike(niko, seattle)
	rater.AddLike(aubreigh, seattle)

	rater.AddLike(niko, denver)
	rater.AddLike(aubreigh, denver)
}
