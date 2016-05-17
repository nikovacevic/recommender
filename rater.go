// rater drives the rating system of the recommender. It manages like/dislike
// data and provides an API for manipulating ratings.

package recommender

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/boltdb/bolt"
)

const ratingsDbName string = "ratings.db"
const userBucketName string = "user"
const itemBucketName string = "item"

type Rater struct {
	db *bolt.DB
}

// NewRater returns a new Rater
func NewRater() *Rater {
	// create key/value store for ratings data
	db, err := bolt.Open(ratingsDbName, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	// create buckets for users and items
	err = db.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists([]byte(userBucketName))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte(itemBucketName))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	return &Rater{db}
}

// Close closes the rater's store connection. Deferring a call to this method
// is recommended on creation of a Rater.
func (r *Rater) Close() {
	err := r.db.Close()
	if err != nil {
		log.Fatal(err)
	}
}

// AddLike records a user liking an item. If the user already likes the item,
// nothing happens. Only if the recording fails will this return an error.
func (r *Rater) AddLike(user *User, item *Item) error {
	var items []Item
	var err error
	// Check if like exists
	err = r.db.View(func(tx *bolt.Tx) error {
		users := tx.Bucket([]byte(userBucketName))
		key := []byte(user.Id)
		if data := users.Get(key); data != nil {
			err = json.Unmarshal(data, &items)
			if err != nil {
				return err
			}
			fmt.Println(items)
		}
		return nil
	})
	if err != nil {
		return err
	}
	err = r.db.Update(func(tx *bolt.Tx) error {
		items = append(items, *item)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

// RemoveLike removes a user's record of liking an item. If the user already
// does not like the item (which is different than disliking it), then
// nothing happens. Only if the removal fails will this return an error.
func (r *Rater) RemoveLike(user *User, item *Item) error {
	return nil
}

func (r *Rater) GetItemsByUser(user *User) ([]Item, error) {
	return make([]Item, 0), nil
}

func (r *Rater) GetUsersByItem(item *Item) ([]User, error) {
	return make([]User, 0), nil
}
