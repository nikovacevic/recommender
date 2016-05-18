// rater drives the rating system of the recommender. It manages like/dislike
// data and provides an API for manipulating ratings.

package recommender

import (
	"bytes"
	"encoding/json"
	"log"
	"os"

	"github.com/boltdb/bolt"
)

const ratingsDbName string = "ratings.db"
const userBucketName string = "user"
const itemBucketName string = "item"

var traceLog *log.Logger

type Rater struct {
	db *bolt.DB
}

func init() {
	traceLog = log.New(os.Stdout, "TRACE: ", log.Ldate|log.Ltime|log.Lshortfile)
}

// NewRater returns a new Rater
func NewRater() *Rater {
	// create key/value store for ratings data
	db, err := bolt.Open(ratingsDbName, 0600, nil)
	if err != nil {
		log.Panic(err)
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
		log.Panic(err)
	}
	return &Rater{db}
}

// Close closes the rater's store connection. Deferring a call to this method
// is recommended on creation of a Rater.
func (r *Rater) Close() {
	err := r.db.Close()
	if err != nil {
		log.Panic(err)
	}
}

// AddLike records a user liking an item. If the user already likes the item,
// nothing happens. Only if the recording fails will this return an error.
func (r *Rater) AddLike(user *User, item *Item) error {
	traceLog.Printf("AddLike (%s, %s)\n", user.Name, item.Name)
	// itemIds used to unmarshal values of key, value pairs, which
	// are (user ID, [item ID, ...])
	var itemIds [][]byte
	// true if (key, value) pair exists in DB
	exists := false

	// TODO If item does not exist in item bucket, add it

	// Return early if user already likes item
	if err := r.db.View(func(tx *bolt.Tx) error {
		userBucket := tx.Bucket([]byte(userBucketName))
		if data := userBucket.Get(user.Id); data != nil {
			// Get user's item IDs
			if err := json.Unmarshal(data, &itemIds); err != nil {
				return err
			}
			// Look for given item ID
			for _, id := range itemIds {
				if bytes.Equal(id, item.Id) {
					exists = true
					break
				}
			}
		}
		return nil
	}); err != nil {
		return err
	}
	if exists {
		traceLog.Printf("Like (%s, %s) already exists.\n", user.Name, item.Name)
		return nil
	}

	// Add item to user's Likes
	if err := r.db.Update(func(tx *bolt.Tx) error {
		userBucket := tx.Bucket([]byte(userBucketName))
		itemIds = append(itemIds, item.Id)
		data, err := json.Marshal(itemIds)
		if err != nil {
			return err
		}
		if err := userBucket.Put(user.Id, data); err != nil {
			return err
		}
		traceLog.Printf("Like (%s, %s) added.\n", user.Name, item.Name)
		return nil
	}); err != nil {
		log.Panic(err)
	}

	// TODO Add user to item's LikedBys

	return nil
}

// RemoveLike removes a user's record of liking an item. If the user already
// does not like the item (which is different than disliking it), then
// nothing happens. Only if the removal fails will this return an error.
func (r *Rater) RemoveLike(user *User, item *Item) error {
	return nil
}

func (r *Rater) GetItemsByUser(user *User) ([]Item, error) {
	// itemIds used to unmarshal values of (key, value) pairs, which
	// are (user ID, [item ID, ...])
	var itemIds [][]byte
	var items []Item
	// Get item IDs
	if err := r.db.View(func(tx *bolt.Tx) error {
		userBucket := tx.Bucket([]byte(userBucketName))
		if data := userBucket.Get(user.Id); data != nil {
			if err := json.Unmarshal(data, &itemIds); err != nil {
				return err
			}
		}
		return nil

	}); err != nil {
		return nil, err
	}
	// Get Items by itemIds
	if err := r.db.View(func(tx *bolt.Tx) error {
		itemBucket := tx.Bucket([]byte(itemBucketName))
		for _, id := range itemIds {
			if data := itemBucket.Get(id); data != nil {
				var item Item
				if err := json.Unmarshal(data, &item); err != nil {
					return err
				}
				items = append(items, item)
			}
		}
		return nil

	}); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *Rater) GetUsersByItem(item *Item) ([]User, error) {
	return make([]User, 0), nil
}
