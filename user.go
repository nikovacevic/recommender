// user defines a user of recommender

package recommender

import (
	"fmt"
	"log"

	uuid "github.com/satori/go.uuid"
)

type User struct {
	Id    []byte `json:"id"`
	Name  string `json:"name"`
	Likes []Item `json:"likes"`
}

func NewUser(name string) *User {
	id, err := uuid.NewV4().MarshalBinary()
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return &User{Id: id, Name: name}
}

func (u *User) String() string {
	return fmt.Sprintf("%s", u.Name)
}
