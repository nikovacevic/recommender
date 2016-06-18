package recommender

import "fmt"

type Score int

const (
	like    Score = 1
	dislike Score = -1
)

type Rating struct {
	Item  Item  `json:"item"`
	Score Score `json:"score"`
}

// String represents an Item as a string
func (r Rating) String() string {
	var score string
	switch r.Score {
	case like:
		score = "like"
	case dislike:
		score = "dislike"
	default:
		score = "none"
	}
	return fmt.Sprintf("%s: %s", r.Item.Name, score)
}
