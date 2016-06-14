package recommender

import (
	"fmt"

	uuid "github.com/satori/go.uuid"
)

// User is ___
type User struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Likes []Item `json:"likes"`
}

// NewUser creates and returns a new User
func NewUser(name string) *User {
	return &User{Id: uuid.NewV4().String(), Name: name}
}

// String represents a User as a string
func (u *User) String() string {
	return fmt.Sprintf("%s", u.Name)
}
