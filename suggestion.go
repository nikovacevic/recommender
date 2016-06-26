package recommender

type SuggestionIndex float32

type Suggestion struct {
	Item  Item            `json:"item"`
	Index SuggestionIndex `json:"index"`
}
