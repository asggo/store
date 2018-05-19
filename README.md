# store
The store package provides a simple wrapper for the Bolt key/value database. store allows you to create and delete buckets in the root of the database and allows you to read, write, and delete key/value pairs within a bucket. Store does not support nested buckets.

    import "github.com/asggo/store/"

## Usage

### Storing Key/Value Pairs
All keys are stored in buckets. Keys can be created and updated with the Write method. Keys can be viewed with the Read method and deleted with the Delete method.

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


### Searching for Buckets and Keys
The AllBuckets and FindBuckets methods allow you to get a list of buckets. The AllKeys and FindKeys methods allow you to get a list of keys from a bucket. The FindBuckets and FindKeys methods will return any items that contain the given search term.

    s := NewStore("/path/to/database/file")

    keys, err := s.AllKeys("bucketname")
    if err != nil {
        fmt.Println("Could not get keys.")
    }

    for _, key := range keys {
        // do something with key
    }

    // Get all keys in the test bucket with key in the name.
    keys, err := s.FindKeys("test", "key")
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

## Methods
Create a new store object with a bolt database located at filePath.

    func NewStore(filePath string) (*Store, error)

AllBuckets returns a list of all the buckets in the root of the database.

    func (s *Store) AllBuckets() ([]string, error)

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

FindBuckets returns all buckets, whose name contains the given string.

    func (s *Store) FindBuckets(needle string) ([]string, error)

FindKeys returns all keys, whose name contains the given string, from the given bucket.

    func (s *Store) FindKeys(bucket, needle string) ([]string, error)

Read gets the value associated with the given key in the given bucket. If the key does not exist, Read returns nil.

    func (s *Store) Read(bucket, key string) []byte

Write stores the given key/value pair in the given bucket.

    func (s *Store) Write(bucket, key string, value []byte) error

# kv
In addition to the store package this repository includes kv.go, which can be compiled and added to the system PATH. kv provides a scriptable interface to a store database.

## Compilation

    go build -o kv kv.go

## Usage

    kv filename action [arguments]

    Actions:
    add <bucketname>                 Add a new bucket to the database.
    add <bucketname> <key> <value>   Add the key/value to the bucket.
    get                              Get a list of buckets.
    get <bucketname>                 Get a list of keys in a bucket.
    get <bucketname> <key>           Get the value of the key in the bucket.
    del <bucketname>                 Delete the bucket and its keys.
    del <bucketname> <key>           Delete the key/value in the bucket
    find <string>                    Find buckets, which contain the string.
    find <bucketname> <string>       Find keys in the bucket, which contain
                                     the string.
