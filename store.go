package store

import (
    "fmt"
    "strings"
    "time"

    "github.com/boltdb/bolt"
)

var BucketNotExist = fmt.Errorf("store: bucket does not exist.")
var BucketNotCreated = fmt.Errorf("store: bucket not created.")

type Store struct {
    db  *bolt.DB
}

func (s *Store) CreateBucket(bucket string) error {
    return s.db.Update(func(tx *bolt.Tx) error {
        _, err := tx.CreateBucketIfNotExists([]byte(bucket))
        if err != nil {
            return BucketNotCreated
        }

        return nil
    })
}

func (s *Store) DeleteBucket(bucket string) error {
    return s.db.Update(func(tx *bolt.Tx) error {
        return tx.DeleteBucket([]byte(bucket))
    })
}

func (s *Store) Create(bucket, key string, value []byte) error {
    err := s.db.Update(func(tx *bolt.Tx) error {
        b := tx.Bucket([]byte(bucket))
        if b == nil {
            return BucketNotExist
        }

        return b.Put([]byte(key), []byte(value))
    })

    return err
}

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

func (s *Store) Update(bucket, key string, value []byte) error {
    err := s.Delete(bucket, key)
    if err != nil {
        return err
    }

    err = s.Create(bucket, key, value)
    if err != nil {
        return err
    }

    return nil
}

func (s *Store) AllKeys(bucket string) ([]string, error) {
    var keys []string

    err := s.db.Update(func(tx *bolt.Tx) error {
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

func (s *Store) FindKeys(bucket, needle string) ([]string, error) {
    var keys []string

    keys, err := s.AllKeys(bucket)
    if err != nil {
        return keys, err
    }

    for _, key := range keys {
        if strings.Contains(key, needle) {
            keys = append(keys, key)
        }
    }

    return keys, nil
}

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
