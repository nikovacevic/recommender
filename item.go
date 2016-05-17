// item defines an item that the recommender can recommend.

package recommender

import (
	"fmt"
	"log"

	uuid "github.com/satori/go.uuid"
)

type Item struct {
	Id      []byte `json:"id"`
	Name    string `json:"name"`
	LikedBy []User `json:"likedBy"`
}

func NewItem(name string) *Item {
	id, err := uuid.NewV4().MarshalBinary()
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return &Item{Id: id, Name: name}
}

func (i *Item) String() string {
	return fmt.Sprintf("%s: %s", i.Id, i.Name)
}
