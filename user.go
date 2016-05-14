// user defines a user of recommender

package recommender

import (
	"fmt"

	uuid "github.com/satori/go.uuid"
)

type User struct {
	Id   uuid.UUID
	Name string
}

func NewUser(name string) *User {
	return &User{Name: name}
}

func (i *User) encode() ([]byte, error) {
	// TODO
	return make([]byte, 0), nil
}

func decode(data []byte) (*User, error) {
	// TODO
	return &User{}
}

func (u *User) String() string {
	return fmt.Sprintf("%s: %s", u.Id.String(), u.Name)
}
