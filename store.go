// store provides a simple wrapper for the Bolt key/value database. store allows
// you to create and delete buckets in the root of the database and allows you
// to read, write, and delete key/value pairs within a bucket. Currently, store
// does not support nested buckets.
//
// Usage:
//
// Storing Key/Value Pairs
// s := NewStore("/path/to/database/file")
// err := s.CreateBucket("bucketname")
// if err != nil {
//     log.Println("Could not create bucket.")
// }
//
// err = s.Write("bucketname", "key", []byte("value"))
// if err != nil {
//     log.Println("Could not write key/value pair.")
// }
//
// val = s.Read("bucketname", "key")
// err = s.Delete("bucketname", "key")
// if err != nil {
//     log.Println("Could not delete key.")
// }
//
//
// Searching for Keys
// s := NewStore("/path/to/database/file")
// keys, err := s.AllKeys("bucketname")
// if err != nil {
//     fmt.Println("Could not get keys.")
// }
//
// for _, key := range keys {
//     // do something with key
// }
//
// // Get all keys with bucket in the name.
// keys, err := s.FindKeys("bucket")
// if err != nil {
//     fmt.Println("Could not get keys.")
// }
//
// for _, key := range keys {
//     // do something with key
// }
package store

import (
	"fmt"
	"strings"
	"time"

	"github.com/boltdb/bolt"
)

var BucketNotExist = fmt.Errorf("store: bucket does not exist.")
var BucketNotCreated = fmt.Errorf("store: bucket not created.")

// Store holds the bolt database
type Store struct {
	db *bolt.DB
}

// CreateBucket creates a new bucket with the given name at the root of the
// database. An error is returned if the bucket cannot be created.
func (s *Store) CreateBucket(bucket string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return BucketNotCreated
		}

		return nil
	})
}

// DeleteBucket deletes the bucket with the given name from the root of the
// database. Returns an error if the bucket cannot be deleted.
func (s *Store) DeleteBucket(bucket string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		return tx.DeleteBucket([]byte(bucket))
	})
}

// AllBuckets returns a list of all the buckets in the root of the database.
func (s *Store) AllBuckets() ([]string, error) {
	var buckets []string

	err := s.db.View(func(tx *bolt.Tx) error {
		tx.ForEach(func(k, v []byte) error {
			buckets = append(buckets, string(k))
			return nil
		})

		return nil
	})

	if err != nil {
		return buckets, err
	}

	return buckets, nil
}

// FindBuckets returns all buckets, whose name contains the given string.
func (s *Store) FindBuckets(needle string) ([]string, error) {
	var buckets []string

	allBuckets, err := s.AllBuckets()
	if err != nil {
		return buckets, err
	}

	for _, bucket := range allBuckets {
		if strings.Contains(bucket, needle) {
			buckets = append(buckets, bucket)
		}
	}

	return buckets, nil
}

// Write stores the given key/value pair in the given bucket.
func (s *Store) Write(bucket, key string, value []byte) error {
	err := s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return BucketNotExist
		}

		return b.Put([]byte(key), []byte(value))
	})

	return err
}

// Read gets the value associated with the given key in the given bucket. If the
// key does not exist, Read returns nil.
func (s *Store) Read(bucket, key string) []byte {
	var val []byte

	s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return BucketNotExist
		}

		val = b.Get([]byte(key))

		return nil
	})

	return val
}

// Delete removes a key/value pair from the given bucket. An error is returned
// if the key/value pair cannot be deleted.
func (s *Store) Delete(bucket, key string) error {
	err := s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return BucketNotExist
		}

		return b.Delete([]byte(key))
	})

	return err
}

// AllKeys returns all of the keys in the given bucket.
func (s *Store) AllKeys(bucket string) ([]string, error) {
	var keys []string

	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return BucketNotExist
		}

		b.ForEach(func(k, v []byte) error {
			keys = append(keys, string(k))
			return nil
		})

		return nil
	})

	if err != nil {
		return keys, err
	}

	return keys, nil
}

// FindKeys returns all keys, whose name contains the given string, from the
// given bucket.
func (s *Store) FindKeys(bucket, needle string) ([]string, error) {
	var keys []string

	allKeys, err := s.AllKeys(bucket)
	if err != nil {
		return keys, err
	}

	for _, key := range allKeys {
		if strings.Contains(key, needle) {
			keys = append(keys, key)
		}
	}

	return keys, nil
}

// Close closes the connection to the bolt database.
func (s *Store) Close() error {
	return s.db.Close()
}

// Create a new store object with a bolt database located at filePath.
func NewStore(filePath string) (*Store, error) {
	s := new(Store)

	db, err := bolt.Open(filePath, 0640, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return s, err
	}

	s.db = db

	return s, nil
}
