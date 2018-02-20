# store
import "github.com/averagesecurityguy/store"

store provides a simple wrapper for the Bolt key/value database. store allows you to create and delete buckets in the root of the database and allows you to read, write, and delete key/value pairs within a bucket. Currently, store does not support nested buckets.

## Usage

### Storing Key/Value Pairs
    s := NewStore("/path/to/database/file")
    err := s.CreateBucket("bucketname")

    if err != nil {
        log.Println("Could not create bucket.")
    }

    err = s.Write("bucketname", "key", []byte("value"))
    if err != nil {
        log.Println("Could not write key/value pair.")
    }

    val = s.Read("bucketname", "key")
    err = s.Delete("bucketname", "key")
    if err != nil {
        log.Println("Could not delete key.")
    }

### Searching for Keys
The AllKeys and FindKeys methods allow you to get a list of keys from a bucket. The FindKeys method will return any keys that contain the given search term.

    s := NewStore("/path/to/database/file")
    keys, err := s.AllKeys("bucketname")
    if err != nil {
        fmt.Println("Could not get keys.")
    }

    for _, key := range keys {
        // do something with key
    }

    // Get all keys with bucket in the name. keys, err :=
    s.FindKeys("bucket")

    if err != nil {
        fmt.Println("Could not get keys.")
    }

    for _, key := range keys {
        // do something with key
    }

## Variables

    var BucketNotCreated = fmt.Errorf("store: bucket not created.")
    var BucketNotExist = fmt.Errorf("store: bucket does not exist.")

## Types
Store holds the bolt database

    type Store struct {
        // contains filtered or unexported fields
    }

Create a new store object with a bolt database located at filePath.

    func NewStore(filePath string) (*Store, error)

AllKeys returns all of the keys in the given bucket.

    func (s *Store) AllKeys(bucket string) ([]string, error)

Close closes the connection to the bolt database.

    func (s *Store) Close() error

CreateBucket creates a new bucket with the given name at the root of the database. An error is returned if the bucket cannot be created.

    func (s *Store) CreateBucket(bucket string) error

Delete removes a key/value pair from the given bucket. An error is returned if the key/value pair cannot be deleted.

    func (s *Store) Delete(bucket, key string) error

DeleteBucket deletes the bucket with the given name from the root of the database. Returns an error if the bucket cannot be deleted.

    func (s *Store) DeleteBucket(bucket string) error

FindKeys returns all keys, whose name contains the given string, from the given bucket.

    func (s *Store) FindKeys(bucket, needle string) ([]string, error)

Read gets the value associated with the given key in the given bucket. If the key does not exist, Read returns nil.

    func (s *Store) Read(bucket, key string) []byte

Write stores the given key/value pair in the given bucket.

    func (s *Store) Write(bucket, key string, value []byte) error
