package recommender

import (
	"bytes"
	"encoding/json"
	"log"
	"os"

	"github.com/boltdb/bolt"
)

const dbName string = "recommender.db"
const userBucketName string = "user"
const userLikesBucketName string = "userLikes"
const itemBucketName string = "item"
const itemLikesBucketName string = "itemLikes"
const userSimilarsBucketName string = "userSimilars"
const itemSimilarsBucketName string = "itemSimilars"

var traceLog *log.Logger

//
type Recommender struct {
	db *bolt.DB
}

func init() {
	traceLog = log.New(os.Stdout, "TRACE: ", log.Ldate|log.Ltime|log.Lshortfile)
}

// NewRecommender returns a new Recommender. The database is opened and buckets are created.
func NewRecommender() (*Recommender, error) {
	// Create key/value store for ratings data
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		return nil, err
	}
	// Create buckets
	if err := db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists([]byte(userBucketName)); err != nil {
			return err
		}
		if _, err := tx.CreateBucketIfNotExists([]byte(itemBucketName)); err != nil {
			return err
		}
		if _, err := tx.CreateBucketIfNotExists([]byte(itemLikesBucketName)); err != nil {
			return err
		}
		if _, err := tx.CreateBucketIfNotExists([]byte(userLikesBucketName)); err != nil {
			return err
		}
		if _, err := tx.CreateBucketIfNotExists([]byte(userSimilarsBucketName)); err != nil {
			return err
		}
		if _, err := tx.CreateBucketIfNotExists([]byte(itemSimilarsBucketName)); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return &Recommender{db}, nil
}

// Close closes the rater's store connection. Deferring a call to this method
// is recommended on creation of a Recommender.
func (r *Recommender) Close() {
	err := r.db.Close()
	if err != nil {
		log.Panic(err)
	}
}

// GetLikesItems gets Items liked by the given User.
func (r *Recommender) GetLikedItems(user *User) ([]Item, error) {
	var items []Item
	var itemIds [][]byte

	if err := r.db.View(func(tx *bolt.Tx) error {
		// Get item IDs
		userLikesBucket := tx.Bucket([]byte(userLikesBucketName))
		data := userLikesBucket.Get(user.Id)
		if data == nil {
			return nil
		}
		if err := json.Unmarshal(data, &itemIds); err != nil {
			return err
		}
		// Get items by ID
		for _, id := range itemIds {
			item, err := r.getItem(id)
			if err != nil {
				traceLog.Printf("WARNING: Cannot find item ID=%v\n", id)
				continue
			}
			items = append(items, *item)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return items, nil
}

// getItem retrieves an Item by ID.
func (r *Recommender) getItem(id []byte) (*Item, error) {
	var item Item

	if err := r.db.View(func(tx *bolt.Tx) error {
		itemBucket := tx.Bucket([]byte(itemBucketName))
		if data := itemBucket.Get(id); data != nil {
			if err := json.Unmarshal(data, &item); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return &item, nil
}

// GetUsersWhoLike retrieves the collection of users who like the given Item.
func (r *Recommender) GetUsersWhoLike(item *Item) ([]User, error) {
	var users []User
	var userIds [][]byte

	if err := r.db.View(func(tx *bolt.Tx) error {
		itemLikeBucket := tx.Bucket([]byte(itemLikesBucketName))
		// Get user IDs
		data := itemLikeBucket.Get(item.Id)
		if data == nil {
			return nil
		}
		if err := json.Unmarshal(data, &userIds); err != nil {
			return err
		}
		// Get users by ID
		for _, id := range userIds {
			user, err := r.getUser(id)
			if err != nil {
				traceLog.Printf("WARNING: Cannot find user ID=%v\n", id)
				continue
			}
			users = append(users, *user)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return users, nil
}

// getUser retrieves a User by ID.
func (r *Recommender) getUser(id []byte) (*User, error) {
	var user User

	if err := r.db.View(func(tx *bolt.Tx) error {
		userBucket := tx.Bucket([]byte(userBucketName))
		if data := userBucket.Get(id); data != nil {
			if err := json.Unmarshal(data, &user); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return &user, nil
}

// AddLike records a user liking an item. If the user already likes the item,
// nothing happens. Only if the recording fails will this return an error.
func (r *Recommender) AddLike(user *User, item *Item) error {
	traceLog.Printf("AddLike (%s, %s)\n", user.Name, item.Name)

	// Add user if record does not already exist
	if err := r.addUser(user); err != nil {
		return err
	}

	// Add item if records does not already exist
	if err := r.addItem(item); err != nil {
		return err
	}

	// Add like (bi-directional) if records do not already exist.
	if err := r.addLike(user, item); err != nil {
		return err
	}

	// In series:
	// TODO go update similarity indices
	// TODO go update suggestions

	return nil
}

// addUser inserts a record in the user bucket if it does not already exist.
func (r *Recommender) addUser(user *User) error {
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

// addItem inserts a record in the item bucket if it does not already exist.
func (r *Recommender) addItem(item *Item) error {
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

// addLike inserts records in the userLikes and itemLikes buckets for the User
// and Item. If either record already exists, no action is takes.
func (r *Recommender) addLike(user *User, item *Item) error {
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
		data, err = json.Marshal(userIds)
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
	return nil
}

// RemoveLike removes a user's record of liking an item. If the user already
// does not like the item (which is different than disliking it), then
// nothing happens. Only if the removal fails will this return an error.
func (r *Recommender) RemoveLike(user *User, item *Item) error {
	// TODO
	return nil
}

// GetItemsByUser retrieves the collection of Items the User likes.
func (r *Recommender) GetItemsByUser(user *User) ([]Item, error) {
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

// GetUsersByItem retrieves the collection of Users who like the given Item.
func (r *Recommender) GetUsersByItem(item *Item) ([]User, error) {
	return make([]User, 0), nil
}

// GetUsers retrieves a collection of Users.
func (r *Recommender) GetUsers(startAt int, count int) ([]User, error) {
	var users []User

	if err := r.db.View(func(tx *bolt.Tx) error {
		userBucket := tx.Bucket([]byte(userBucketName))
		cur := userBucket.Cursor()
		i, c := 0, 0
		for key, val := cur.First(); key != nil && c < count; key, val = cur.Next() {
			if i >= startAt {
				var u User
				if err := json.Unmarshal(val, &u); err != nil {
					return err
				}
				users = append(users, u)
				c++
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return users, nil
}

// GetItems retrieves a collection of Items.
func (r *Recommender) GetItems(startAt int, count int) ([]Item, error) {
	var items []Item

	if err := r.db.View(func(tx *bolt.Tx) error {
		itemBucket := tx.Bucket([]byte(itemBucketName))
		cur := itemBucket.Cursor()
		i, c := 0, 0
		for key, val := cur.First(); key != nil && c < count; key, val = cur.Next() {
			if i > startAt {
				var i Item
				if err := json.Unmarshal(val, &i); err != nil {
					return err
				}
				items = append(items, i)
				c++
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return items, nil
}
