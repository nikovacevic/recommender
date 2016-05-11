package recommender

type User struct {
	Id   int
	Name string
}

func NewUser(name string) *User {
	return &User{Name: name}
}
