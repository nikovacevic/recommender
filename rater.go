// rater drives the rating system of the recommender. It manages like/dislike
// data and provides an API for manipulating ratings.
//
// Data structure
//
// User bucket
// user: [
//	id: {
//		name: "Name",
//		like: [
//			item,
//			...
//		],
//		dislike: [
//			item,
//			...
//		]
//	},
//	...
// ]
//
// Item bucket
// item: [
//	id: {
//		name: "Name",
//		like: [
//			user,
//			...
//		],
//		dislike: [
//			user,
//			...
//		]
//	},
//	...
// ]
//
//

package recommender

import (
	"log"

	"github.com/boltdb/bolt"
)

const ratingsDb string = "ratings.db"
const userBucket []byte = []byte("user")
const itemBucket []byte = []byte("item")

type Rater struct {
	db *bolt.DB
}

// NewRater returns a new Rater
func NewRater() *Rater {
	// create key/value store for ratings data
	db, err := bolt.Open(ratingsDb, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	// create buckets for likes and dislikes
	err := db.Update(func(tx *bolt.Tx) error {
		like, err := tx.CreateBucketIfNotExists(likeBucket)
		if err != nil {
			return err
		}
		dislike, err := tx.CreateBucketIfNotExists(dislikeBucket)
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
func (r *Rater) AddLike(user User, item Item) error {
	var items []Items
	// Check if like exists
	err := r.db.View(func(tx *bolt.Tx) error {
		user := tx.Bucket(userBucket)
		key := []byte(user.Id)
		if val := likes.Get(key); val != nil {
			items = json.
		}
	})
	if err != nil {
		return err
	}
	err := r.db.Update(func(tx *bolt.Tx) error {

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
func (r *Rater) RemoveLike(user User, item Item) error {
	return nil
}

// AddDisike records a user disliking an item. If the user already dislikes the
// item, nothing happens. Only if the recording fails will this return an error.
func (r *Rater) AddDislike(user User, item Item) error {
	err := r.db.Update(func(tx *bolt.Tx) error {

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

// RemoveDislike removes a user's record of disliking an item. If the user
// already does not dislike the item (which is different than liking it), then
// nothing happens. Only if the removal fails will this return an error.
func (r *Rater) RemoveDislike(user User, item Item) error {
	return nil
}

func (r *Rater) GetItemsByUser(user User) []Item, error {
	return make([]Item, 0), nil
}

func (r *Rater) GetUsersByItem(item Item) []User, error {
	return make([]User, 0), nil
}
