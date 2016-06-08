package recommender_test

import (
	"log"
	"testing"

	"github.com/nikovacevic/recommender"
)

const (
	pass = "\u2713"
	fail = "\u2717"
)

func TestAddLike(t *testing.T) {
	rater, err := recommender.NewRater()
	if err != nil {
		log.Fatal(err)
	}
	defer rater.Close()

	niko := recommender.NewUser("Niko Kovacevic")

	denver := recommender.NewItem("Denver")
	phoenix := recommender.NewItem("Phoenix")
	pittsburgh := recommender.NewItem("Pittsburgh")
	portland := recommender.NewItem("Portland")

	items, err := rater.GetLikedItems(niko)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	if l := len(items); l != 0 {
		t.Errorf("There should be zero items. There are %d.", l)
	}

	rater.AddLike(niko, phoenix)
	rater.AddLike(niko, denver)
	rater.AddLike(niko, pittsburgh)

	items, err = rater.GetLikedItems(niko)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	if l := len(items); l != 3 {
		t.Errorf("There should be three items. There are %d.", l)
	}

	rater.AddLike(niko, phoenix)
	rater.AddLike(niko, portland)
	rater.AddLike(niko, pittsburgh)

	items, err = rater.GetLikedItems(niko)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	if l := len(items); l != 4 {
		t.Errorf("There should be four items. There are %d.", l)
	}
}
