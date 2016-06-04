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
const userLikesBucketName string = "userLikes"
const itemBucketName string = "item"
const itemLikesBucketName string = "itemLikes"

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
		_, err = tx.CreateBucketIfNotExists([]byte(itemLikesBucketName))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte(userLikesBucketName))
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

	// If user doesn't exist, add to user bucket
	if err := r.addUser(user); err != nil {
		return err
	}

	// If item doesn't exist, add to item bucket
	if err := r.addItem(item); err != nil {
		return err
	}

	// TODO Make private method
	// Add item to user's likes
	if err := r.db.Update(func(tx *bolt.Tx) error {
		itemIds := [][]byte{}
		userLikesBucket := tx.Bucket([]byte(userLikesBucketName))
		// Find user's liked items
		if data := userLikesBucket.Get(user.Id); data != nil {
			// Get user's liked item IDs
			if err := json.Unmarshal(data, &itemIds); err != nil {
				return err
			}
			// Look for given item ID
			for _, id := range itemIds {
				if bytes.Equal(id, item.Id) {
					// If user already likes item, return early
					traceLog.Printf("Like by (%s, %s) already exists.\n", item.Name, user.Name)
					return nil
				}
			}
		}
		// Add item to user's liked items
		itemIds = append(itemIds, item.Id)
		data, err := json.Marshal(itemIds)
		if err != nil {
			return err
		}
		if err := userLikesBucket.Put(user.Id, data); err != nil {
			return err
		}
		traceLog.Printf("Like by (%s, %s) added.\n", item.Name, user.Name)
		return nil
	}); err != nil {
		return err
	}

	// TODO Return early if the user already likes the item?

	// TODO Make private method
	// Add user to item's LikedBys
	if err := r.db.Update(func(tx *bolt.Tx) error {
		userIds := [][]byte{}
		itemLikesBucket := tx.Bucket([]byte(itemLikesBucketName))
		// Find users who like item
		if data := itemLikesBucket.Get(item.Id); data != nil {
			// Get item's liked-by user IDs
			if err := json.Unmarshal(data, &userIds); err != nil {
				return err
			}
			// Look for user ID
			for _, id := range userIds {
				if bytes.Equal(id, user.Id) {
					// If item already liked by user, return early
					traceLog.Printf("Like (%s, %s) already exists.\n", user.Name, item.Name)
					return nil
				}
			}
		}
		// Add user to item's liked-by
		userIds = append(userIds, user.Id)
		data, err := json.Marshal(userIds)
		if err != nil {
			return err
		}
		if err := itemLikesBucket.Put(item.Id, data); err != nil {
			return err
		}
		traceLog.Printf("Like (%s, %s) added.\n", user.Name, item.Name)
		return nil
	}); err != nil {
		return err
	}

	// In series, do the following (first should `go call` the next)
	// TODO go routine: update similarity indices
	// TODO go routine: update suggestions

	return nil
}

func (r *Rater) addUser(user *User) error {
	err := r.db.Update(func(tx *bolt.Tx) error {
		userBucket := tx.Bucket([]byte(userBucketName))
		// Return early if user already exists
		if data := userBucket.Get(user.Id); data != nil {
			return nil
		}
		// Write JSON encoding of user to DB
		data, err := json.Marshal(user)
		if err != nil {
			return err
		}
		if err := userBucket.Put(user.Id, data); err != nil {
			return err
		}
		traceLog.Printf("User %s added.\n", user.Name)
		return nil
	})

	if err != nil {
		return err
	}
	return nil
}

func (r *Rater) addItem(item *Item) error {
	err := r.db.Update(func(tx *bolt.Tx) error {
		itemBucket := tx.Bucket([]byte(itemBucketName))
		// Return early if item already exists
		if data := itemBucket.Get(item.Id); data != nil {
			return nil
		}
		// Write JSON encoding of item to DB
		data, err := json.Marshal(item)
		if err != nil {
			return err
		}
		if err := itemBucket.Put(item.Id, data); err != nil {
			return err
		}
		traceLog.Printf("Item %s added.\n", item.Name)
		return nil
	})

	if err != nil {
		return err
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
