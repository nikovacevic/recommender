package recommender

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/boltdb/bolt"
)

type Recommender struct {
	db *bolt.DB
}

const (
	dbName                   string = "recommender.db"
	userBucketName           string = "user"
	itemBucketName           string = "item"
	userLikesBucketName      string = "userLikes"
	itemLikesBucketName      string = "itemLikes"
	userDislikesBucketName   string = "userDislikes"
	itemDislikesBucketName   string = "itemDislikes"
	userSimilarityBucketName string = "userSimilarity"
	itemSimilarityBucketName string = "itemSimilarity"
	suggestionBucketName     string = "suggestionBucket"
)

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
		if _, err := tx.CreateBucketIfNotExists([]byte(itemDislikesBucketName)); err != nil {
			return err
		}
		if _, err := tx.CreateBucketIfNotExists([]byte(userDislikesBucketName)); err != nil {
			return err
		}
		if _, err := tx.CreateBucketIfNotExists([]byte(userSimilarityBucketName)); err != nil {
			return err
		}
		if _, err := tx.CreateBucketIfNotExists([]byte(itemSimilarityBucketName)); err != nil {
			return err
		}
		if _, err := tx.CreateBucketIfNotExists([]byte(suggestionBucketName)); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return &Recommender{db}, nil
}

// Close closes the Recommender's store connection. Deferring a call to this method
// is recommended on creation of a Recommender.
func (r *Recommender) Close() {
	err := r.db.Close()
	if err != nil {
		log.Panic(err)
	}
}

// GetLikedItems gets Items liked by the given User.
func (r *Recommender) GetLikedItems(user *User) (map[string]Item, error) {
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
				log.Printf("WARNING: Cannot find item ID=%v\n", id)
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
func (r *Recommender) GetDislikedItems(user *User) (map[string]Item, error) {
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
				log.Printf("WARNING: Cannot find item ID=%v\n", id)
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
func (r *Recommender) getItem(id string) (*Item, error) {
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
func (r *Recommender) GetUsersWhoLike(item *Item) (map[string]User, error) {
	users := make(map[string]User)
	userIds := make(map[string]bool)

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
				log.Printf("WARNING: Cannot find user ID=%v\n", id)
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
func (r *Recommender) GetUsersWhoDislike(item *Item) (map[string]User, error) {
	users := make(map[string]User)
	userIds := make(map[string]bool)

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
				log.Printf("WARNING: Cannot find user ID=%v\n", id)
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
func (r *Recommender) GetUsersWhoRated(item *Item) (map[string]User, error) {
	userCh := make(chan User)
	var wg sync.WaitGroup

	// Retrieve users who like the item and pipe them into the channel.
	wg.Add(1)
	go func() {
		users, err := r.GetUsersWhoLike(item)
		if err != nil {
			return
		}
		for _, user := range users {
			userCh <- user
		}
		wg.Done()
	}()

	// Retrieve users who dislike the item and pipe them into the channel.
	wg.Add(1)
	go func() {
		users, err := r.GetUsersWhoDislike(item)
		if err != nil {
			return
		}
		for _, user := range users {
			userCh <- user
		}
		wg.Done()
	}()

	// Wait for the like and dislike goroutines to finish, then close the
	// user channel
	go func() {
		wg.Wait()
		close(userCh)
	}()

	// As users are sent through the channel, build out map. Return map when
	// channel closes.
	users := make(map[string]User)
	for user := range userCh {
		users[user.Id] = user
	}

	return users, nil
}

// getUser retrieves a User by ID.
func (r *Recommender) getUser(id string) (*User, error) {
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
func (r *Recommender) Like(user *User, item *Item) error {
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

	// Update similarity index
	if err := r.UpdateSimilarity(user); err != nil {
		return err
	}

	// Update suggestions
	// TODO
	if err := r.UpdateSuggestions(user); err != nil {
		return err
	}

	return nil
}

// Dislike records a user disliking an item. If the user already dislikes the
// item, nothing happens. If the user likes the item, the like is removed first.
// Only if the recording fails will this return an error.
func (r *Recommender) Dislike(user *User, item *Item) error {
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

	// Update similarity index
	if err := r.UpdateSimilarity(user); err != nil {
		return err
	}

	// Update suggestions
	// TODO
	if err := r.UpdateSuggestions(user); err != nil {
		return err
	}

	return nil
}

// addUser inserts a record in the user bucket if it does not already exist.
func (r *Recommender) addUser(user *User) error {
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
		// log.Printf("User %s added.\n", user.Name)
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
		// log.Printf("Item %s added.\n", item.Name)
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
func (r *Recommender) addLike(user *User, item *Item) error {
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
				// log.Printf("Like by (%s, %s) already exists.\n", item.Name, user.Name)
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
				// log.Printf("Disike (%s, %s) removed.\n", user.Name, item.Name)
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
				// log.Printf("Like (%s, %s) already exists.\n", user.Name, item.Name)
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
		userIds = make(map[string]bool)
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
				// log.Printf("Dislike (%s, %s) removed.\n", user.Name, item.Name)
			}
		}

		// log.Printf("Like (%s, %s) added.\n", user.Name, item.Name)
		return nil
	}); err != nil {
		return err
	}
	return nil
}

// addDislike inserts records in the userDislikes and itemDislikes buckets for
// the User and Item. If a like exists, both such records are deleted. If the
// dislike records already exists, no action is taken.
func (r *Recommender) addDislike(user *User, item *Item) error {
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
				// log.Printf("Dislike (%s, %s) already exists.\n", item.Name, user.Name)
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
				// log.Printf("Like (%s, %s) removed.\n", user.Name, item.Name)
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
				// log.Printf("Dislike (%s, %s) already exists.\n", user.Name, item.Name)
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
		userIds = make(map[string]bool)
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
				// log.Printf("Like (%s, %s) removed.\n", user.Name, item.Name)
			}
		}

		// log.Printf("Dislike (%s, %s) added.\n", user.Name, item.Name)
		return nil
	}); err != nil {
		return err
	}
	return nil
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

// channelRatings retrieves a collection of Items.
func (r *Recommender) channelRatings(user *User) (<-chan Rating, error) {
	ratingCh := make(chan Rating)
	var wg sync.WaitGroup

	// Retrieve liked items, package them into Rating structs, and pipe
	// them into the rating channel.
	wg.Add(1)
	go func() {
		defer wg.Done()
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
	}()

	// Retrieve disliked items, package them into Rating structs, and pipe
	// them into the rating channel.
	wg.Add(1)
	go func() {
		defer wg.Done()
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
	}()

	// Wait for the like and dislike goroutines to finish, then close the
	// rating channel
	go func() {
		wg.Wait()
		close(ratingCh)
	}()

	return ratingCh, nil
}

// GetRatings retrieves all items a user has rated and returns a map of
// item ID to Rating, which includes the item and the score the user gave.
func (r *Recommender) GetRatings(user *User) (map[string]Rating, error) {
	// As ratings are sent through the rating channel, build out rating
	// map. Return map when channel closes.
	ratings := make(map[string]Rating)
	ratingCh, err := r.channelRatings(user)
	if err != nil {
		return nil, err
	}
	for rating := range ratingCh {
		ratings[rating.Item.Id] = rating
	}

	return ratings, nil
}

// GetRatingNeighbors returns a set of users, indexed by user ID, who rated the
// same items that the given user rated.
func (r *Recommender) GetRatingNeighbors(user *User) (map[string]User, error) {
	neighborMap := make(map[string]User)

	// Get user's ratings
	// TODO If user's ratings are already populated, skip this?
	// user.Ratings == nil did not work as intended
	ratings, err := r.GetRatings(user)
	if err != nil {
		return nil, err
	}
	user.Ratings = ratings

	// Add to neighborMap users who have also rated each item
	for _, rating := range user.Ratings {
		item := rating.Item
		neighbors, err := r.GetUsersWhoRated(&item)
		if err != nil {
			return nil, err
		}
		for id, neighbor := range neighbors {
			// Skip neighbors that have already been added
			if _, exists := neighborMap[id]; !exists {
				// Get the neighbor's ratings
				neighborRatings, err := r.GetRatings(&neighbor)
				if err != nil {
					return nil, err
				}
				neighbor.Ratings = neighborRatings
				// Add neighbor to map
				neighborMap[id] = neighbor
			}
		}
	}

	// Delete given user from their own set of neighbors
	delete(neighborMap, user.Id)

	return neighborMap, nil
}

// UpdateSimilarity calculates the similarity index for each user with which the
// given user has overlapping rated items.
func (r *Recommender) UpdateSimilarity(user *User) error {
	// Get user's rated items
	// TODO If user's ratings are already populated, skip this?
	// user.Ratings == nil did not work as intended
	ratings, err := r.GetRatings(user)
	if err != nil {
		return err
	}
	user.Ratings = ratings

	// Get user's neighbors
	neighbors, err := r.GetRatingNeighbors(user)
	if err != nil {
		return err
	}

	// Compute similarity index for each of user's neighbors
	// Run each neighbor concurrently, but wait for completion of all
	var wg sync.WaitGroup
	similarityCh := make(chan *Similarity)
	for _, neighbor := range neighbors {
		wg.Add(1)
		// Create new instance of neighbor for goroutine
		neighbor := neighbor
		go func() {
			index := r.similarityIndex(user, &neighbor)
			similarityCh <- &Similarity{
				User:  neighbor,
				Index: index,
			}
			wg.Done()
		}()
	}

	// Close similarity channel when all goroutines complete
	go func() {
		wg.Wait()
		close(similarityCh)
	}()

	// Map neighbor's user ID to similarity index
	for similarity := range similarityCh {
		// Update database
		r.updateSimilarity(user, &similarity.User, similarity.Index)
	}

	return nil
}

// similarityIndex calculates the current similarity index based on the ratings
// in each user's Ratings
func (r *Recommender) similarityIndex(user1, user2 *User) SimilarityIndex {
	var agree, disagree int

	for id, rating1 := range user1.Ratings {
		if rating2, exists := user2.Ratings[id]; exists {
			if rating1.Score == rating2.Score {
				agree++
			} else {
				disagree++
			}
		}
	}

	index := float32((agree - disagree)) / float32((agree + disagree))

	return SimilarityIndex(index)
}

// updateSimilarity updates the similarity index for the given users
func (r *Recommender) updateSimilarity(user1 *User, user2 *User, index SimilarityIndex) error {
	if err := r.db.Update(func(tx *bolt.Tx) error {
		userSimilarityBucket := tx.Bucket([]byte(userSimilarityBucketName))

		// Get user1's existing similarities
		similarityMap := make(map[string]SimilarityIndex)
		if data := userSimilarityBucket.Get([]byte(user1.Id)); data != nil {
			if err := json.Unmarshal(data, &similarityMap); err != nil {
				return err
			}
		}
		// Set new index
		similarityMap[user2.Id] = index
		data, err := json.Marshal(similarityMap)
		if err != nil {
			return err
		}
		// Write updated map to bucket
		if err := userSimilarityBucket.Put([]byte(user1.Id), data); err != nil {
			return err
		}

		// Get user2's existing similarities
		similarityMap = make(map[string]SimilarityIndex)
		if data := userSimilarityBucket.Get([]byte(user2.Id)); data != nil {
			if err := json.Unmarshal(data, &similarityMap); err != nil {
				return err
			}
		}
		// Set new index
		similarityMap[user1.Id] = index
		data, err = json.Marshal(similarityMap)
		if err != nil {
			return err
		}
		// Write updated map to bucket
		if err := userSimilarityBucket.Put([]byte(user2.Id), data); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

// channelSimilarity returns a channel of the given user's similarities
func (r *Recommender) channelSimilarity(user *User) (<-chan Similarity, error) {
	similarityCh := make(chan Similarity)
	similarityIndexMap := make(map[string]SimilarityIndex)

	if err := r.db.View(func(tx *bolt.Tx) error {
		userSimilarityBucket := tx.Bucket([]byte(userSimilarityBucketName))
		if data := userSimilarityBucket.Get([]byte(user.Id)); data != nil {
			if err := json.Unmarshal(data, &similarityIndexMap); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}

	go func() {
		for id, index := range similarityIndexMap {
			u, err := r.getUser(id)
			if err != nil {
				close(similarityCh)
				return
			}
			similarityCh <- Similarity{
				User:  *u,
				Index: index,
			}
		}
		close(similarityCh)
	}()

	return similarityCh, nil
}

// GetSimilarity returns a map the given user's similarities, keyed by their
// similar user's ID
func (r *Recommender) GetSimilarity(user *User) (map[string]Similarity, error) {
	similarityMap := make(map[string]Similarity)
	similarityCh, err := r.channelSimilarity(user)
	if err != nil {
		return nil, err
	}
	for similarity := range similarityCh {
		similarityMap[similarity.User.Id] = similarity
	}
	return similarityMap, nil
}

// UpdateSuggestions generates a set of Suggestions (items with corresponding
// suggestion index) for the given user.
func (r *Recommender) UpdateSuggestions(user *User) error {
	//log.Printf("UpdateSuggestions(%s)\n", user.Name)

	// Get similarities for user
	similarityMap, err := r.GetSimilarity(user)
	if err != nil {
		return err
	}

	// For each similarity, get similar user's rated items, but only items
	// user has not rated.
	itemCh := make(chan Item)
	var wg sync.WaitGroup
	for _, similarity := range similarityMap {
		wg.Add(1)
		// Create new instance of similarity for goroutine
		similarity := similarity
		go func() {
			// Get similar user's ratings
			ratingsCh, err := r.channelRatings(&(similarity.User))
			if err != nil {
				return
			}
			// For each rated item, if user has not rated the item,
			// send it into itemCh
			for r := range ratingsCh {
				if _, exists := user.Ratings[r.Item.Id]; !exists {
					itemCh <- r.Item
				}
			}
			wg.Done()
		}()
	}

	go func() {
		defer close(itemCh)
		wg.Wait()
	}()

	// For each item, suggestion index = (zL-zD)/total, where zL is the sum
	// of the similarity indices of users who like the item, zD is the sum
	// of the similarity indices of users who dislike the item, and total is
	// the total number of users composing zL and zD.
	// TODO Can we optimize this?
	suggestionMap := make(map[string]Suggestion)
	for item := range itemCh {
		// Get all users who have rated the item, separated into like
		// and dislike.
		likeUsers, err := r.GetUsersWhoLike(&item)
		if err != nil {
			return err
		}
		dislikeUsers, err := r.GetUsersWhoDislike(&item)
		if err != nil {
			return err
		}
		// Scan each similar user for a like/dislike score. If one
		// exists, increment index parameters.
		var zL, zD, total float32
		for id, similarity := range similarityMap {
			if _, exists := likeUsers[id]; exists {
				zL += float32(similarity.Index)
				total++
			} else if _, exists := dislikeUsers[id]; exists {
				zD += float32(similarity.Index)
				total++
			}
		}
		// Build Suggestion, then add it to the map
		index := (zL - zD) / total
		suggestionMap[item.Id] = Suggestion{
			Item:  item,
			Index: SuggestionIndex(index),
		}
	}

	// Save the suggestion map, keyed by the user's Id
	if err := r.db.Update(func(tx *bolt.Tx) error {
		suggestionBucket := tx.Bucket([]byte(suggestionBucketName))
		data, err := json.Marshal(suggestionMap)
		if err != nil {
			return err
		}
		if err := suggestionBucket.Put([]byte(user.Id), data); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}

// GetSuggestions retrieves the set of Suggestions for the given user.
func (r *Recommender) GetSuggestions(user *User) (map[string]Suggestion, error) {
	suggestionMap := make(map[string]Suggestion)
	if err := r.db.View(func(tx *bolt.Tx) error {
		suggestionBucket := tx.Bucket([]byte(suggestionBucketName))
		if data := suggestionBucket.Get([]byte(user.Id)); data != nil {
			if err := json.Unmarshal(data, &suggestionMap); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return suggestionMap, nil
}
