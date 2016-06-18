package recommender

import (
	"fmt"

	uuid "github.com/satori/go.uuid"
)

type Item struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

// NewItem creates and returns an Item
func NewItem(name string) *Item {
	return &Item{Id: uuid.NewV4().String(), Name: name}
}

// String represents an Item as a string
func (i Item) String() string {
	return fmt.Sprintf("%s", i.Name)
}
