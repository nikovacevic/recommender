package recommender

import (
	"fmt"
	"log"

	uuid "github.com/satori/go.uuid"
)

// Item is ___
type Item struct {
	Id      []byte `json:"id"`
	Name    string `json:"name"`
	LikedBy []User `json:"likedBy"`
}

// NewItem creates and returns an Item
func NewItem(name string) *Item {
	id, err := uuid.NewV4().MarshalBinary()
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return &Item{Id: id, Name: name}
}

// String represents an Item as a string
func (i *Item) String() string {
	return fmt.Sprintf("%s", i.Name)
}
