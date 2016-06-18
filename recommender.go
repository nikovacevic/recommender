package recommender

import (
	"encoding/json"
	"log"
	"os"
	"sync"

	"github.com/boltdb/bolt"
)

//
type Rater struct {
	db *bolt.DB
}

const (
	dbName                 string = "Rater.db"
	userBucketName         string = "user"
	itemBucketName         string = "item"
	userLikesBucketName    string = "userLikes"
	itemLikesBucketName    string = "itemLikes"
	userDislikesBucketName string = "userDislikes"
	itemDislikesBucketName string = "itemDislikes"
	userSimilarsBucketName string = "userSimilars"
	itemSimilarsBucketName string = "itemSimilars"
)

var traceLog *log.Logger

func init() {
	traceLog = log.New(os.Stdout, "TRACE: ", log.Ldate|log.Ltime|log.Lshortfile)
}

// NewRater returns a new Rater. The database is opened and buckets are created.
func NewRater() (*Rater, error) {
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
		if _, err := tx.CreateBucketIfNotExists([]byte(itemDislikesBucketName)); err != nil {
			return err
		}
		if _, err := tx.CreateBucketIfNotExists([]byte(userDislikesBucketName)); err != nil {
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
	return &Rater{db}, nil
}

// Close closes the Rater's store connection. Deferring a call to this method
// is recommended on creation of a Rater.
func (r *Rater) Close() {
	err := r.db.Close()
	if err != nil {
		log.Panic(err)
	}
}

// GetLikedItems gets Items liked by the given User.
func (r *Rater) GetLikedItems(user *User) (map[string]Item, error) {
	items := make(map[string]Item)
	itemIds := make(map[string]bool)

	if err := r.db.View(func(tx *bolt.Tx) error {
		// Get item IDs
		userLikesBucket := tx.Bucket([]byte(userLikesBucketName))
		data := userLikesBucket.Get([]byte(user.Id))
		if data == nil {
			return nil
		}
		if err := json.Unmarshal(data, &itemIds); err != nil {
			return err
		}
		// Get items by ID
		for id, _ := range itemIds {
			item, err := r.getItem(id)
			if err != nil {
				traceLog.Printf("WARNING: Cannot find item ID=%v\n", id)
				continue
			}
			items[id] = *item
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return items, nil
}

// GetDislikedItems gets Items disliked by the given User.
func (r *Rater) GetDislikedItems(user *User) (map[string]Item, error) {
	items := make(map[string]Item)
	itemIds := make(map[string]bool)

	if err := r.db.View(func(tx *bolt.Tx) error {
		// Get item IDs
		userDislikesBucket := tx.Bucket([]byte(userDislikesBucketName))
		data := userDislikesBucket.Get([]byte(user.Id))
		if data == nil {
			return nil
		}
		if err := json.Unmarshal(data, &itemIds); err != nil {
			return err
		}
		// Get items by ID
		for id, _ := range itemIds {
			item, err := r.getItem(id)
			if err != nil {
				traceLog.Printf("WARNING: Cannot find item ID=%v\n", id)
				continue
			}
			items[id] = *item
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return items, nil
}

// getItem retrieves an Item by ID.
func (r *Rater) getItem(id string) (*Item, error) {
	var item Item

	if err := r.db.View(func(tx *bolt.Tx) error {
		itemBucket := tx.Bucket([]byte(itemBucketName))
		if data := itemBucket.Get([]byte(id)); data != nil {
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
func (r *Rater) GetUsersWhoLike(item *Item) (map[string]User, error) {
	var users map[string]User
	var userIds map[string]bool

	if err := r.db.View(func(tx *bolt.Tx) error {
		itemLikeBucket := tx.Bucket([]byte(itemLikesBucketName))
		// Get user IDs
		data := itemLikeBucket.Get([]byte(item.Id))
		if data == nil {
			return nil
		}
		if err := json.Unmarshal(data, &userIds); err != nil {
			return err
		}
		// Get users by ID
		for id, _ := range userIds {
			user, err := r.getUser(id)
			if err != nil {
				traceLog.Printf("WARNING: Cannot find user ID=%v\n", id)
				continue
			}
			users[id] = *user
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return users, nil
}

// GetUsersWhoDislike retrieves the collection of users who dislike the given Item.
func (r *Rater) GetUsersWhoDislike(item *Item) (map[string]User, error) {
	var users map[string]User
	var userIds map[string]bool

	if err := r.db.View(func(tx *bolt.Tx) error {
		itemDislikeBucket := tx.Bucket([]byte(itemDislikesBucketName))
		// Get user IDs
		data := itemDislikeBucket.Get([]byte(item.Id))
		if data == nil {
			return nil
		}
		if err := json.Unmarshal(data, &userIds); err != nil {
			return err
		}
		// Get users by ID
		for id, _ := range userIds {
			user, err := r.getUser(id)
			if err != nil {
				traceLog.Printf("WARNING: Cannot find user ID=%v\n", id)
				continue
			}
			users[id] = *user
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return users, nil
}

// GetUsersWhoRated retrieves the collection of users who rated the given Item.
func (r *Rater) GetUsersWhoRated(item *Item) (map[string]User, error) {
	var users map[string]User
	var userIds map[string]bool

	if err := r.db.View(func(tx *bolt.Tx) error {
		itemDislikeBucket := tx.Bucket([]byte(itemDislikesBucketName))
		// Get user IDs
		data := itemDislikeBucket.Get([]byte(item.Id))
		if data == nil {
			return nil
		}
		if err := json.Unmarshal(data, &userIds); err != nil {
			return err
		}
		// Get users by ID
		for id, _ := range userIds {
			user, err := r.getUser(id)
			if err != nil {
				traceLog.Printf("WARNING: Cannot find user ID=%v\n", id)
				continue
			}
			users[id] = *user
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return users, nil
}

// getUser retrieves a User by ID.
func (r *Rater) getUser(id string) (*User, error) {
	var user User

	if err := r.db.View(func(tx *bolt.Tx) error {
		userBucket := tx.Bucket([]byte(userBucketName))
		if data := userBucket.Get([]byte(id)); data != nil {
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

// Like records a user liking an item. If the user already likes the item,
// nothing happens. Only if the recording fails will this return an error.
func (r *Rater) Like(user *User, item *Item) error {
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

// Dislike records a user disliking an item. If the user already dislikes the
// item, nothing happens. If the user likes the item, the like is removed first.
// Only if the recording fails will this return an error.
func (r *Rater) Dislike(user *User, item *Item) error {
	// Add user if record does not already exist
	if err := r.addUser(user); err != nil {
		return err
	}

	// Add item if records does not already exist
	if err := r.addItem(item); err != nil {
		return err
	}

	// Add like (bi-directional) if records do not already exist.
	if err := r.addDislike(user, item); err != nil {
		return err
	}

	// In series:
	// TODO go update similarity indices
	// TODO go update suggestions

	return nil
}

// addUser inserts a record in the user bucket if it does not already exist.
func (r *Rater) addUser(user *User) error {
	err := r.db.Update(func(tx *bolt.Tx) error {
		userBucket := tx.Bucket([]byte(userBucketName))
		// Return early if user already exists
		if data := userBucket.Get([]byte(user.Id)); data != nil {
			return nil
		}
		// Write JSON encoding of user to DB
		data, err := json.Marshal(user)
		if err != nil {
			return err
		}
		if err := userBucket.Put([]byte(user.Id), data); err != nil {
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
func (r *Rater) addItem(item *Item) error {
	err := r.db.Update(func(tx *bolt.Tx) error {
		itemBucket := tx.Bucket([]byte(itemBucketName))
		// Return early if item already exists
		if data := itemBucket.Get([]byte(item.Id)); data != nil {
			return nil
		}
		// Write JSON encoding of item to DB
		data, err := json.Marshal(item)
		if err != nil {
			return err
		}
		if err := itemBucket.Put([]byte(item.Id), data); err != nil {
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
// and Item. If a dislike exists, both such records are deleted. If the like
// records already exists, no action is taken.
func (r *Rater) addLike(user *User, item *Item) error {
	// Add item to user's likes
	if err := r.db.Update(func(tx *bolt.Tx) error {
		itemIds := make(map[string]bool)
		userLikesBucket := tx.Bucket([]byte(userLikesBucketName))

		// Find user's liked items
		if data := userLikesBucket.Get([]byte(user.Id)); data != nil {
			// Get user's liked item IDs
			if err := json.Unmarshal(data, &itemIds); err != nil {
				return err
			}
			// If user already likes item, return early
			if itemIds[item.Id] {
				traceLog.Printf("Like by (%s, %s) already exists.\n", item.Name, user.Name)
				return nil
			}
		}

		// Add item to user's liked items
		itemIds[item.Id] = true
		data, err := json.Marshal(itemIds)
		if err != nil {
			return err
		}
		if err := userLikesBucket.Put([]byte(user.Id), data); err != nil {
			return err
		}

		// Find user's disliked items
		itemIds = make(map[string]bool)
		userDislikesBucket := tx.Bucket([]byte(userDislikesBucketName))
		if data := userDislikesBucket.Get([]byte(user.Id)); data != nil {
			// Get user's disliked item IDs
			if err := json.Unmarshal(data, &itemIds); err != nil {
				return err
			}
			// If user dislikes item, remove the dislike
			if itemIds[item.Id] {
				delete(itemIds, item.Id)
				data, err := json.Marshal(itemIds)
				if err != nil {
					return err
				}
				if err := userDislikesBucket.Put([]byte(user.Id), data); err != nil {
					return err
				}
			}
		}

		userIds := make(map[string]bool)
		itemLikesBucket := tx.Bucket([]byte(itemLikesBucketName))
		// Find users who like item
		if data := itemLikesBucket.Get([]byte(item.Id)); data != nil {
			// Get item's liked-by user IDs
			if err := json.Unmarshal(data, &userIds); err != nil {
				return err
			}
			// If item already liked by user, return early
			if userIds[user.Id] {
				traceLog.Printf("Like (%s, %s) already exists.\n", user.Name, item.Name)
				return nil
			}
		}

		// Add user to item's liked-by
		userIds[user.Id] = true
		data, err = json.Marshal(userIds)
		if err != nil {
			return err
		}
		if err := itemLikesBucket.Put([]byte(item.Id), data); err != nil {
			return err
		}

		// Find users who dislike item
		itemDislikesBucket := tx.Bucket([]byte(itemDislikesBucketName))
		if data := itemDislikesBucket.Get([]byte(item.Id)); data != nil {
			// Get item's liked item IDs
			if err := json.Unmarshal(data, &userIds); err != nil {
				return err
			}
			// If item dislikes item, remove the dislike
			if userIds[user.Id] {
				delete(userIds, user.Id)
				data, err := json.Marshal(userIds)
				if err != nil {
					return err
				}
				if err := itemDislikesBucket.Put([]byte(user.Id), data); err != nil {
					return err
				}
				traceLog.Printf("Dislike (%s, %s) removed.\n", user.Name, item.Name)
			}
		}

		traceLog.Printf("Like (%s, %s) added.\n", user.Name, item.Name)
		return nil
	}); err != nil {
		return err
	}
	return nil
}

// addDislike inserts records in the userDislikes and itemDislikes buckets for
// the User and Item. If a like exists, both such records are deleted. If the
// dislike records already exists, no action is taken.
func (r *Rater) addDislike(user *User, item *Item) error {
	// Add item to user's dislikes
	if err := r.db.Update(func(tx *bolt.Tx) error {
		itemIds := make(map[string]bool)
		userDislikesBucket := tx.Bucket([]byte(userDislikesBucketName))

		// Find user's disliked items
		if data := userDislikesBucket.Get([]byte(user.Id)); data != nil {
			// Get user's disliked item IDs
			if err := json.Unmarshal(data, &itemIds); err != nil {
				return err
			}
			// If user already dislikes item, return early
			if itemIds[item.Id] {
				traceLog.Printf("Dislike (%s, %s) already exists.\n", item.Name, user.Name)
				return nil
			}
		}

		// Add item to user's disliked items
		itemIds[item.Id] = true
		data, err := json.Marshal(itemIds)
		if err != nil {
			return err
		}
		if err := userDislikesBucket.Put([]byte(user.Id), data); err != nil {
			return err
		}

		// Find user's liked items
		itemIds = make(map[string]bool)
		userLikesBucket := tx.Bucket([]byte(userLikesBucketName))
		if data := userLikesBucket.Get([]byte(user.Id)); data != nil {
			// Get user's liked item IDs
			if err := json.Unmarshal(data, &itemIds); err != nil {
				return err
			}
			// If user likes item, remove the like
			if itemIds[item.Id] {
				delete(itemIds, item.Id)
				data, err := json.Marshal(itemIds)
				if err != nil {
					return err
				}
				if err := userLikesBucket.Put([]byte(user.Id), data); err != nil {
					return err
				}
				traceLog.Printf("Like (%s, %s) removed.\n", user.Name, item.Name)
			}
		}

		userIds := make(map[string]bool)
		itemDislikesBucket := tx.Bucket([]byte(itemDislikesBucketName))
		// Find users who dislike item
		if data := itemDislikesBucket.Get([]byte(item.Id)); data != nil {
			// Get item's disliked-by user IDs
			if err := json.Unmarshal(data, &userIds); err != nil {
				return err
			}
			// If item already disliked by user, return early
			if userIds[user.Id] {
				traceLog.Printf("Dislike (%s, %s) already exists.\n", user.Name, item.Name)
				return nil
			}
		}

		// Add user to item's disliked-by
		userIds[user.Id] = true
		data, err = json.Marshal(userIds)
		if err != nil {
			return err
		}
		if err := itemDislikesBucket.Put([]byte(item.Id), data); err != nil {
			return err
		}

		// Find users who like item
		itemLikesBucket := tx.Bucket([]byte(itemLikesBucketName))
		if data := itemLikesBucket.Get([]byte(item.Id)); data != nil {
			// Get item's liked item IDs
			if err := json.Unmarshal(data, &userIds); err != nil {
				return err
			}
			// If item dislikes item, remove the dislike
			if userIds[user.Id] {
				delete(userIds, user.Id)
				data, err := json.Marshal(userIds)
				if err != nil {
					return err
				}
				if err := itemLikesBucket.Put([]byte(user.Id), data); err != nil {
					return err
				}
			}
		}

		traceLog.Printf("Dislike (%s, %s) added.\n", user.Name, item.Name)
		return nil
	}); err != nil {
		return err
	}
	return nil
}

// GetUsers retrieves a collection of Users.
func (r *Rater) GetUsers(startAt int, count int) ([]User, error) {
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
func (r *Rater) GetItems(startAt int, count int) ([]Item, error) {
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

// GetRatings retrieves all items a user has rated and returns a map of
// item ID to Rating, which includes the item and the score the user gave.
func (r *Rater) GetRatings(user *User) (map[string]Rating, error) {
	ratingCh := make(chan Rating)
	var wg sync.WaitGroup

	// Retrieve liked items, package them into Rating structs, and pipe
	// them into the rating channel.
	wg.Add(1)
	go func() {
		items, err := r.GetLikedItems(user)
		if err != nil {
			return
		}
		for _, item := range items {
			ratingCh <- Rating{
				Item:  item,
				Score: like,
			}
		}
		wg.Done()
	}()

	// Retrieve disliked items, package them into Rating structs, and pipe
	// them into the rating channel.
	wg.Add(1)
	go func() {
		items, err := r.GetDislikedItems(user)
		if err != nil {
			return
		}
		for _, item := range items {
			ratingCh <- Rating{
				Item:  item,
				Score: dislike,
			}
		}
		wg.Done()
	}()

	// Wait for the like and dislike goroutines to finish, then close the
	// rating channel
	go func() {
		wg.Wait()
		close(ratingCh)
	}()

	// As ratings are sent through the rating channel, build out rating
	// map. Return map when channel closes.
	ratings := make(map[string]Rating)
	for rating := range ratingCh {
		ratings[rating.Item.Id] = rating
	}

	return ratings, nil
}
