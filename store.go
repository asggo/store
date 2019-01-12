// store provides a simple wrapper for the bbolt key/value database. store
// allows you to create and delete buckets in the root of the database and
// allows you to read, write, and delete key/value pairs within a bucket.
// Currently, store does not support nested buckets.
package store

import (
	"bytes"
	"fmt"
	"os"
	"time"

	bolt "go.etcd.io/bbolt"
)

// WalkFunc is called for each key/value pair when walking the database.
type WalkFunc func(key string, val []byte)

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
			return fmt.Errorf("store: bucket %s not created: %s", bucket, err)
		}

		return nil
	})
}

// DeleteBucket deletes the bucket with the given name from the root of the
// database. Returns an error if the bucket cannot be deleted.
func (s *Store) DeleteBucket(bucket string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket([]byte(bucket))

		if err != nil {
			return fmt.Errorf("store: could not delete bucket %s: %s", bucket, err)
		}

		return nil
	})
}

// Walk executes the WalkFunc on each bucket in the root.
func (s *Store) Walk(fn WalkFunc) error {
	return s.db.View(func(tx *bolt.Tx) error {
		c := tx.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			fn(string(k), v)
		}

		return nil
	})
}

// WalkBucket executes the WalkBucketFunc on each key, value pair in the bucket.
func (s *Store) WalkBucket(bucket string, fn WalkFunc) error {
	return s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return fmt.Errorf("store: bucket %s does not exist", bucket)
		}

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			fn(string(k), v)
		}

		return nil
	})
}

// WalkPrefix executes the WalkFunc on every key/value pair in a bucket where
// the key matches the given prefix.
func (s *Store) WalkPrefix(bucket, prefix string, fn WalkFunc) error {
	return s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return fmt.Errorf("store: bucket %s does not exist", bucket)
		}

		c := b.Cursor()
		pre := []byte(prefix)

		for k, v := c.Seek(pre); k != nil && bytes.HasPrefix(k, pre); k, v = c.Next() {
			fn(string(k), v)
		}

		return nil
	})
}

// Read key/value pairs from a bucket in batches of count size. Update the
// batch with the found items. On error, the key/value map will be nil and
// should not be used.
func (s *Store) ReadBatch(bucket, next string, count int) (map[string][]byte, string, error) {
	var items map[string][]byte

	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return fmt.Errorf("store: bucket %s does not exist", bucket)
		}

		items = make(map[string][]byte)
		c := b.Cursor()

		for k, v := c.Seek([]byte(next)); k != nil && len(items) < count; k, v = c.Next() {
			items[string(k)] = v
			next = string(k)
		}

		if len(items) != count {
			next = ""
		}

		return nil
	})

	if err != nil {
		return nil, "", err
	}

	return items, next, nil
}

// Write stores the given key/value pair in the given bucket.
func (s *Store) Write(bucket, key string, value []byte) error {
	return s.db.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return fmt.Errorf("store: bucket %s does not exist", bucket)
		}

		err := b.Put([]byte(key), value)
		if err != nil {
			return fmt.Errorf("store: could not write to key %s in bucket %s: %s", key, bucket, err)
		}

		return nil
	})
}

// Read gets the value associated with the given key in the given bucket. If the
// key does not exist, Read returns nil.
func (s *Store) Read(bucket, key string) ([]byte, error) {
	var val []byte

	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return fmt.Errorf("store: bucket %s does not exist", bucket)
		}

		val = b.Get([]byte(key))
		if val == nil {
			return fmt.Errorf("store: key %s does not exit", key)
		}

		return nil
	})

	return val, err
}

// Delete removes a key/value pair from the given bucket. An error is returned
// if the key/value pair cannot be deleted.
func (s *Store) Delete(bucket, key string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return fmt.Errorf("store: bucket %s does not exist", bucket)
		}

		err := b.Delete([]byte(key))
		if err != nil {
			return fmt.Errorf("store: could not delete key %s in bucket %s", key, bucket)
		}

		return nil
	})
}

// Backup the database to the given file.
func (s *Store) Backup(filename string) error {
	return s.db.View(func(tx *bolt.Tx) error {
		file, err := os.Create(filename)
		if err != nil {
			return fmt.Errorf("store: could not create backup file %s: %s", filename, err)
		}

		defer file.Close()

		_, err = tx.WriteTo(file)
		if err != nil {
			return fmt.Errorf("store: could not write to backup file %s: %s", filename, err)
		}

		return nil
	})
}

// Close closes the connection to the bolt database.
func (s *Store) Close() error {
	err := s.db.Close()
	if err != nil {
		return fmt.Errorf("store: could not close database")
	}

	return nil
}

// Create a new store object with a bolt database located at filePath.
func NewStore(filePath string) (*Store, error) {
	var err error

	s := new(Store)

	for tries := 1; tries < 20; tries += 2 {
		timeout := 1 << uint(tries) * time.Millisecond

		db, err := bolt.Open(filePath, 0640, &bolt.Options{Timeout: timeout})
		if err == nil {
			s.db = db
			return s, nil
		}
	}

	return nil, fmt.Errorf("store: can not open database %s: %s", filePath, err)
}
