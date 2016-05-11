package recommender

import (
	"log"

	"github.com/boltdb/bolt"
)

type RaterKind string

type Rater struct {
	db   *bolt.DB
	kind RaterKind
}

func NewRater(kind string) *Rater {
	// store for like data
	db, err := bolt.Open(kind+".db", 0600, nil)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	return &Rater{db, RaterKind(kind)}
}

func (r *Rater) CloseDB() {
	r.db.Close()
}

func (r *Rater) Add(user User, item Item) {
	err := r.db.Update(func(tx *bolt.Tx) error {
		// TODO
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}

func (r *Rater) Remove(user User, item Item) {

}

func (r *Rater) ItemsByUser(user User) {

}

func (r *Rater) UsersByItem(item Item) {

}
