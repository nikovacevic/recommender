package recommender

type SimilarityIndex float32

type Similarity struct {
	User  User            `json:"user"`
	Index SimilarityIndex `json:"index"`
}
