package recommender

import (
	"fmt"
	"log"

	uuid "github.com/satori/go.uuid"
)

// User is ___
type User struct {
	Id    []byte `json:"id"`
	Name  string `json:"name"`
	Likes []Item `json:"likes"`
}

// NewUser creates and returns a new User
func NewUser(name string) *User {
	id, err := uuid.NewV4().MarshalBinary()
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return &User{Id: id, Name: name}
}

// String represents a User as a string
func (u *User) String() string {
	return fmt.Sprintf("%s", u.Name)
}
