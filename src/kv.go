package main

import (
	"fmt"
	"os"

	"github.com/asggo/store"
)

func help() {
	u := `Usage:
	kv filename action [arguments]

Actions:
	add <bucketname>                 Add a new bucket to the database.
	add <bucketname> <key> <value>   Add the key/value to the bucket.
	get                              Get a list of buckets.
	get <bucketname>                 Get a list of keys in a bucket.
	get <bucketname> <key>           Get the value of the key in the bucket.
	val <bucketname>                 Get a list of values in a bucket.
	val <bucketname> <string>        Get a list of values in a bucket, which
	                                 contain the string.
	del <bucketname>                 Delete the bucket and its keys.
	del <bucketname> <key>           Delete the key/value in the bucket
	find <string>                    Find buckets, which contain the string.
	find <bucketname> <string>       Find keys in the bucket, which contain
	                                 the string.
	backup <filename>                Backup the database to this file.
	`
	fmt.Println(u)
	os.Exit(1)
}

// add <bucketname>                 Adds a new bucket to the database.
// add <bucketname> <key> <value>   Add the key/value to the bucket.
func add(db *store.Store, args []string) {
	switch len(args) {
	case 1:
		err := db.CreateBucket(args[0])
		if err != nil {
			fmt.Printf("Could not create bucket %s: %s\n", args[0], err)
		}
	case 3:
		err := db.Write(args[0], args[1], []byte(args[2]))
		if err != nil {
			fmt.Printf("Could not write to bucket %s: %s\n", args[0], err)
		}
	default:
		help()
	}
}

// get                      Returns a list of buckets.
// get <bucketname>         Returns all keys in a bucket.
// get <bucketname> <key>   Returns the value of the key in the bucket.
func get(db *store.Store, args []string) {
	var items []string
	var err error

	switch len(args) {
	case 0:
		items, err = db.AllBuckets()
		if err != nil {
			fmt.Printf("Could not retrieve bucket list: %s\n", err)
		}
	case 1:
		items, err = db.AllKeys(args[0])
		if err != nil {
			fmt.Printf("Could not retrieve keys from bucket %s: %s\n", args[0], err)
		}
	case 2:
		value := db.Read(args[0], args[1])
		fmt.Println(string(value))
	default:
		help()
	}

	for _, item := range items {
		fmt.Println(item)
	}
}

// delete <bucketname>         Delete the bucket and its keys.
// delete <bucketname> <key>   Delete the key/value in the bucket
func delete(db *store.Store, args []string) {
	switch len(args) {
	case 1:
		err := db.DeleteBucket(args[0])
		if err != nil {
			fmt.Printf("Could not delete bucket %s: %s\n", args[0], err)
		}
	case 2:
		err := db.Delete(args[0], args[1])
		if err != nil {
			fmt.Printf("Could not delete key %s from bucket %s: %s\n", args[0], args[1], err)
		}
	default:
		help()
	}
}

// find <string>                Find all buckets in the database, which contain the string.
// find <bucketname> <string>   Find all keys in the bucket, which contain the string.
func find(db *store.Store, args []string) {
	var items []string
	var err error

	switch len(args) {
	case 1:
		items, err = db.FindBuckets(args[0])
		if err != nil {
			fmt.Printf("Could not find buckets matching %s: %s\n", args[0], err)
		}
	case 2:
		items, err = db.FindKeys(args[0], args[1])
		if err != nil {
			fmt.Printf("Could not find keys matching %s in bucket %s: %s\n", args[1], args[0], err)
			return
		}
	default:
		help()
	}

	for _, item := range items {
		fmt.Println(item)
	}
}

// val <bucketname>            Return all values in the bucket.
// val <bucketname> <string>   Return all values in the bucket, which contain the string.
func val(db *store.Store, args []string) {
	var items []string
	var err error

	switch len(args) {
	case 1:
		items, err = db.AllVals(args[0])
		if err != nil {
			fmt.Printf("Could not get values from bucket %s: %s\n", args[0], err)
		}
	case 2:
		items, err = db.FindVals(args[0], args[1])
		if err != nil {
			fmt.Printf("Could not find values matching %s in bucket %s: %s\n", args[1], args[0], err)
			return
		}
	default:
		help()
	}

	for _, item := range items {
		fmt.Println(item)
	}
}

func backup(db *store.Store, args []string) {
	switch len(args) {
	case 1:
		err := db.Backup(args[0])
		if err != nil {
			fmt.Printf("Could not backup database to %s: %s\n", args[0], err)
			return
		}
	default:
		help()
	}
}

func main() {
	if len(os.Args) < 3 {
		help()
	}

	// Open our database file.
	dbfile := os.Args[1]
	db, err := store.NewStore(dbfile)
	if err != nil {
		fmt.Println("Could not open database file:", err)
	}

	action := os.Args[2]

	switch action {
	case "add":
		add(db, os.Args[3:])
	case "get":
		get(db, os.Args[3:])
	case "val":
		val(db, os.Args[3:])
	case "del":
		delete(db, os.Args[3:])
	case "find":
		find(db, os.Args[3:])
	case "backup":
		backup(db, os.Args[3:])
	default:
		help()
	}

	db.Close()
}
