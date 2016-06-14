package recommender_test

import (
	"log"
	"testing"

	"github.com/nikovacevic/recommender"
)

func TestAddLike(t *testing.T) {
	rater, err := recommender.NewRecommender()
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

	// GetLikedItems should return no items at this point
	items, err := rater.GetLikedItems(niko)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	if l := len(items); l != 0 {
		t.Errorf("There should be zero items. There are %d.", l)
	}

	// GetUsersWhoLike should return no users at this point
	users, err := rater.GetUsersWhoLike(portland)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	if l := len(users); l != 0 {
		t.Errorf("There should be zero users. There are %d.", l)
	}

	// Add some likes
	rater.AddLike(niko, phoenix)
	rater.AddLike(niko, denver)
	rater.AddLike(niko, pittsburgh)
	rater.AddLike(aubreigh, phoenix)
	rater.AddLike(aubreigh, portland)

	// Get the liked items
	items, err = rater.GetLikedItems(niko)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	if l := len(items); l != 3 {
		t.Errorf("There should be three items. There are %d.", l)
	}

	// Add some more likes, with a few overlapping
	rater.AddLike(niko, phoenix)
	rater.AddLike(niko, portland)
	rater.AddLike(niko, pittsburgh)
	rater.AddLike(aubreigh, phoenix)

	// There should only be one new item, four total
	items, err = rater.GetLikedItems(niko)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	if l := len(items); l != 4 {
		t.Errorf("There should be four items. There are %d.", l)
	}
}
