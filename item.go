// item defines an item that the recommender can recommend.

package recommender

import (
	"fmt"

	uuid "github.com/satori/go.uuid"
)

type Item struct {
	Id   uuid.UUID
	Name string
}

func NewItem(name string) *Item {
	return &Item{Name: name}
}

func (i *Item) encode() ([]byte, error) {
	// TODO
	return make([]byte, 0), nil
}

func decode(data []byte) (*Item, error) {
	// TODO
	return &Item{}
}

func (i *Item) String() string {
	return fmt.Sprintf("%s: %s", i.Id.String(), name)
}
