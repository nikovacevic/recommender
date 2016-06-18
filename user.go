package recommender

import (
	"fmt"

	uuid "github.com/satori/go.uuid"
)

type User struct {
	Id      string            `json:"id"`
	Name    string            `json:"name"`
	Ratings map[string]Rating `json:"ratings"`
}

// NewUser creates and returns a new User
func NewUser(name string) *User {
	return &User{Id: uuid.NewV4().String(), Name: name}
}

// String represents a User as a string
func (u User) String() string {
	return fmt.Sprintf("%s", u.Name)
}
