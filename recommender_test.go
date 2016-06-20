package recommender_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/nikovacevic/recommender"
)

func TestLike(t *testing.T) {
	// log.Printf("TestLike")

	rater, err := recommender.NewRater()
	if err != nil {
		log.Fatal(err)
	}
	defer rater.Close()

	niko := recommender.NewUser("Niko Kovacevic")
	aubreigh := recommender.NewUser("Aubreigh Brunschwig")

	denver := recommender.NewItem("Denver")
	phoenix := recommender.NewItem("Phoenix")
	pittsburgh := recommender.NewItem("Pittsburgh")
	portland := recommender.NewItem("Portland")
	losAngeles := recommender.NewItem("Los Angeles")
	miami := recommender.NewItem("Miami")

	// GetLikedItems should return no items at this point
	items, err := rater.GetLikedItems(niko)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	if l := len(items); l != 0 {
		t.Errorf("There should be 0 items. There are %d.", l)
	}

	// GetUsersWhoLike should return no users at this point
	users, err := rater.GetUsersWhoLike(portland)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	if l := len(users); l != 0 {
		t.Errorf("There should be 0 users. There are %d.", l)
	}

	// Add some likes
	rater.Like(niko, phoenix)
	rater.Like(niko, denver)
	rater.Like(niko, pittsburgh)
	rater.Like(aubreigh, phoenix)
	rater.Like(aubreigh, portland)

	// Get the liked items
	items, err = rater.GetLikedItems(niko)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	if l := len(items); l != 3 {
		t.Errorf("There should be 3 items. There are %d.", l)
	}

	// Add some dislikes
	rater.Dislike(niko, phoenix)
	rater.Dislike(niko, miami)
	rater.Dislike(niko, losAngeles)

	// Add some more likes, with some overlapping and previously disliked
	rater.Like(niko, phoenix)
	rater.Like(niko, portland)
	rater.Like(niko, pittsburgh)
	rater.Like(aubreigh, phoenix)

	// There should only be one new item, 4 total
	items, err = rater.GetLikedItems(niko)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	if l := len(items); l != 4 {
		t.Errorf("There should be 4 items. There are %d: %v", l, items)
	}
}

func TestDisLike(t *testing.T) {
	// log.Printf("TestDislike")

	rater, err := recommender.NewRater()
	if err != nil {
		log.Fatal(err)
	}
	defer rater.Close()

	niko := recommender.NewUser("Niko Kovacevic")

	phoenix := recommender.NewItem("Phoenix")
	losAngeles := recommender.NewItem("Los Angeles")
	miami := recommender.NewItem("Miami")

	// GetLikedItems should return no items at this point
	items, err := rater.GetDislikedItems(niko)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	if l := len(items); l != 0 {
		t.Errorf("There should be 0 items. There are %d.", l)
	}

	// GetUsersWhoDislike should return no users at this point
	users, err := rater.GetUsersWhoDislike(miami)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	if l := len(users); l != 0 {
		t.Errorf("There should be 0 users. There are %d.", l)
	}

	// Add some dislikes
	rater.Dislike(niko, phoenix)
	rater.Dislike(niko, miami)

	// Get the disliked items
	items, err = rater.GetDislikedItems(niko)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	if l := len(items); l != 2 {
		t.Errorf("There should be 2 items. There are %d.", l)
	}

	// Like some items
	rater.Like(niko, phoenix)

	// Add some more dislikes, with some overlapping and previously liked
	rater.Dislike(niko, phoenix)
	rater.Dislike(niko, losAngeles)

	// There should only be one new item, 3 total
	items, err = rater.GetDislikedItems(niko)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	if l := len(items); l != 3 {
		t.Errorf("There should be 3 items. There are %d.", l)
	}
}

func TestGetRatings(t *testing.T) {
	// log.Printf("TestGetRatings")

	rater, err := recommender.NewRater()
	if err != nil {
		log.Fatal(err)
	}
	defer rater.Close()

	niko := recommender.NewUser("Niko Kovacevic")

	phoenix := recommender.NewItem("Phoenix")
	losAngeles := recommender.NewItem("Los Angeles")
	miami := recommender.NewItem("Miami")
	pittsburgh := recommender.NewItem("Pittsburgh")
	boulder := recommender.NewItem("Boulder")
	seattle := recommender.NewItem("Seattle")

	// GetLikedItems should return no items at this point
	ratings, err := rater.GetRatings(niko)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	if len(ratings) != 0 {
		t.Errorf("There should be 0 items. There are %d.", len(ratings))
	}

	// Add some likes and dislikes
	rater.Dislike(niko, phoenix)
	rater.Dislike(niko, miami)
	rater.Dislike(niko, losAngeles)
	rater.Like(niko, pittsburgh)
	rater.Like(niko, boulder)
	rater.Like(niko, seattle)

	// There should be six ratings
	ratings, err = rater.GetRatings(niko)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	if len(ratings) != 6 {
		t.Errorf("There should be six items. There are %d.", len(ratings))
	}

	/*
		// Print ratings (notice order varies because of concurrency)
		for _, rating := range ratings {
			fmt.Printf("%v\n", rating)
		}
	*/
}

func TestGetUsersWhoRated(t *testing.T) {
	// log.Printf("TestGetUsersWhoRated")

	rater, err := recommender.NewRater()
	if err != nil {
		log.Fatal(err)
	}
	defer rater.Close()

	niko := recommender.NewUser("Niko Kovacevic")
	aubreigh := recommender.NewUser("Aubreigh Brunschwig")
	johnny := recommender.NewUser("Johnny Bernard")
	amanda := recommender.NewUser("Amanda Hunt")
	nick := recommender.NewUser("Nick Evers")

	phoenix := recommender.NewItem("Phoenix")
	pittsburgh := recommender.NewItem("Pittsburgh")

	// GetUsersWhoRated should return no users at this point
	users, err := rater.GetUsersWhoRated(pittsburgh)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	if len(users) != 0 {
		t.Errorf("There should be 0 users. There are %d.", len(users))
	}

	// Add some likes and dislikes
	rater.Dislike(niko, phoenix)
	rater.Dislike(aubreigh, phoenix)
	rater.Like(johnny, phoenix)
	rater.Like(amanda, phoenix)
	rater.Like(niko, pittsburgh)
	rater.Like(nick, pittsburgh)

	// GetUsersWhoRated should return 4 users at this point
	users, err = rater.GetUsersWhoRated(phoenix)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	if len(users) != 4 {
		t.Errorf("There should be 4 users. There are %d:", len(users))
		for _, user := range users {
			fmt.Printf("%v\n", user)
		}
	}

	// GetUsersWhoRated should return 2 users at this point
	users, err = rater.GetUsersWhoRated(pittsburgh)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	if len(users) != 2 {
		t.Errorf("There should be 2 users. There are %d:", len(users))
		for _, user := range users {
			fmt.Printf("%v\n", user)
		}
	}
}

func TestGetRatingNeighbors(t *testing.T) {
	// log.Printf("TestGetRatingNeighbors")

	rater, err := recommender.NewRater()
	if err != nil {
		log.Fatal(err)
	}
	defer rater.Close()

	niko := recommender.NewUser("Niko Kovacevic")
	aubreigh := recommender.NewUser("Aubreigh Brunschwig")
	johnny := recommender.NewUser("Johnny Bernard")
	amanda := recommender.NewUser("Amanda Hunt")
	nick := recommender.NewUser("Nick Evers")

	phoenix := recommender.NewItem("Phoenix")
	pittsburgh := recommender.NewItem("Pittsburgh")

	// GetRatingNeighbors should return no users at this point
	neighbors, err := rater.GetRatingNeighbors(niko)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	if len(neighbors) != 0 {
		t.Errorf("There should be 0 users. There are %d:", len(neighbors))
		for _, neighbor := range neighbors {
			fmt.Printf("%v\n", neighbor)
		}
	}

	// Add some likes and dislikes
	rater.Dislike(niko, phoenix)
	rater.Dislike(aubreigh, phoenix)
	rater.Like(johnny, phoenix)
	rater.Like(amanda, phoenix)
	rater.Like(niko, pittsburgh)
	rater.Dislike(aubreigh, pittsburgh)
	rater.Like(nick, pittsburgh)

	// GetRatingNeighbors should return five users at this point
	neighbors, err = rater.GetRatingNeighbors(niko)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	if len(neighbors) != 4 {
		t.Errorf("There should be 4 users. There are %d:", len(neighbors))
		for _, neighbor := range neighbors {
			fmt.Printf("%v\n", neighbor)
		}
	}
}

func TestSimilarity(t *testing.T) {
	log.Printf("TestSimilarity")

	rater, err := recommender.NewRater()
	if err != nil {
		log.Fatal(err)
	}
	defer rater.Close()

	niko := recommender.NewUser("Niko Kovacevic")
	aubreigh := recommender.NewUser("Aubreigh Brunschwig")
	johnny := recommender.NewUser("Johnny Bernard")
	nick := recommender.NewUser("Nick Evers")

	phoenix := recommender.NewItem("Phoenix")
	pittsburgh := recommender.NewItem("Pittsburgh")
	boulder := recommender.NewItem("Boulder")
	losAngeles := recommender.NewItem("Los Angeles")
	portland := recommender.NewItem("Portland")
	seattle := recommender.NewItem("Seattle")

	// GetSimilarity should return nothing at this point
	nikoSims, err := rater.GetSimilarity(niko)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	if len(nikoSims) != 0 {
		t.Errorf("There should be 0 similarities. There are %d:", len(nikoSims))
		for _, similarity := range nikoSims {
			fmt.Printf("%v\n", similarity)
		}
	}

	// Add some likes and dislikes
	rater.Dislike(niko, phoenix)
	rater.Like(niko, pittsburgh)
	rater.Like(niko, boulder)
	rater.Dislike(niko, losAngeles)
	rater.Like(niko, portland)
	rater.Like(niko, seattle)

	rater.Dislike(aubreigh, phoenix)
	rater.Dislike(aubreigh, pittsburgh)
	rater.Like(aubreigh, boulder)
	rater.Like(aubreigh, losAngeles)
	rater.Like(aubreigh, portland)
	rater.Like(aubreigh, seattle)

	rater.Like(johnny, phoenix)
	rater.Like(johnny, losAngeles)

	rater.Like(nick, pittsburgh)
	rater.Like(nick, portland)

	// GetSimilarity should return three similarities at this point
	nikoSims, err = rater.GetSimilarity(niko)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	if len(nikoSims) != 3 {
		t.Errorf("There should be 3 similarities. There are %d:", len(nikoSims))
		for _, similarity := range nikoSims {
			fmt.Printf("%v\n", similarity)
		}
	}

	// Get other users's similarities
	aubreighSims, err := rater.GetSimilarity(aubreigh)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	nickSims, err := rater.GetSimilarity(nick)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	johnnySims, err := rater.GetSimilarity(johnny)
	if err != nil {
		t.Errorf("Error: %s", err)
	}

	// Test commutativity
	if nikoSims[aubreigh.Id].Index != aubreighSims[niko.Id].Index {
		t.Errorf("Similarity(Niko, Aubreigh) should equaul Similarity(Aubreigh, Niko).")
	}
	if nikoSims[nick.Id].Index != nickSims[niko.Id].Index {
		t.Errorf("Similarity(Niko, Nick) should equaul Similarity(Nick, Niko).")
	}
	if nikoSims[johnny.Id].Index != johnnySims[niko.Id].Index {
		t.Errorf("Similarity(Niko, Johnny) should equaul Similarity(Johnny, Niko).")
	}

	// Test values
	if float32(nikoSims[aubreigh.Id].Index) != float32(2.0/6.0) {
		t.Errorf("Similarity(Niko, Aubreigh) should be %f. Actually %f", float32(2.0/6.0), nikoSims[aubreigh.Id])
	}
	if float32(nikoSims[nick.Id].Index) != float32(1) {
		t.Errorf("Similarity(Niko, Nick) should be %f. Actually %f", 1, nikoSims[nick.Id])
	}
	if float32(nikoSims[johnny.Id].Index) != float32(-1) {
		t.Errorf("Similarity(Niko, Johnny) should be %f. Actually %f", -1, nikoSims[johnny.Id])
	}
}
