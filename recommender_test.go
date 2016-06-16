package recommender_test

import (
	"log"
	"testing"

	"github.com/nikovacevic/recommender"
)

func TestAddLike(t *testing.T) {
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
		t.Errorf("There should be three items. There are %d.", l)
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

	// There should only be one new item, four total
	items, err = rater.GetLikedItems(niko)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	if l := len(items); l != 4 {
		t.Errorf("There should be four items. There are %d: %v", l, items)
	}
}

func TestDisLike(t *testing.T) {
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
		t.Errorf("There should be zero items. There are %d.", l)
	}

	// GetUsersWhoDislike should return no users at this point
	users, err := rater.GetUsersWhoDislike(miami)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	if l := len(users); l != 0 {
		t.Errorf("There should be zero users. There are %d.", l)
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
		t.Errorf("There should be two items. There are %d.", l)
	}

	// Like some items
	rater.Like(niko, phoenix)

	// Add some more dislikes, with some overlapping and previously liked
	rater.Dislike(niko, phoenix)
	rater.Dislike(niko, losAngeles)

	// There should only be one new item, three total
	items, err = rater.GetDislikedItems(niko)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	if l := len(items); l != 3 {
		t.Errorf("There should be three items. There are %d.", l)
	}
}
