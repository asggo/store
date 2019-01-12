package store

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestStore(t *testing.T) {
	testNewStore(t)
	testBucket(t)
	testKey(t)
	testWalk(t)
}

func testNewStore(t *testing.T) {
	// Open a database in a path that does not exist.
	_, err := NewStore("bad/path/test.db")
	if err == nil {
		t.Fatal("Create Store: expected error got nil")
	}

	// Open a database
	s, err := NewStore("test.db")
	if err != nil {
		t.Fatal("Create Store: unexpected error", err)
	}
	s.Close()
	os.Remove("test.db")
}

func testBucket(t *testing.T) {
	s, _ := NewStore("test.db")
	defer s.Close()
	defer os.Remove("test.db")

	err := s.CreateBucket("bucket")
	if err != nil {
		t.Fatal("Create Bucket: unexpected error", err)
	}

	err = s.DeleteBucket("bucket")
	if err != nil {
		t.Fatal("Delete Bucket: unexpected error", err)
	}

	err = s.DeleteBucket("nonexistent")
	if err == nil {
		t.Fatal("Delete Bucket: expected an error while deleting bucket, got nil")
	}
}

func testKey(t *testing.T) {
	var data1 = []byte("Store this.")
	var data2 = []byte("Store that.")

	// Open a database
	s, _ := NewStore("test.db")
	defer s.Close()
	defer os.Remove("test.db")

	// Create a bucket for storing keys.
	s.CreateBucket("bucket1")

	// Create keys
	err := s.Write("bucket1", "key1", data1)
	if err != nil {
		t.Fatal("Write Key: unexpected error", err)
	}

	s.Write("bucket1", "key2", data2)

	// Read keys
	val, err := s.Read("nonexistent", "key")
	if err == nil {
		t.Error("Read Key: expected error reading from nonexistent bucket")
	}

	val, err = s.Read("bucket1", "nonexistent")
	if err == nil {
		t.Error("Read Key: expected error when reading non-existent key")
	}

	val, _ = s.Read("bucket1", "key1")
	if string(val) != string(data1) {
		t.Error("Read: expected", string(data1), "got", string(val))
	}

	val, _ = s.Read("bucket1", "key2")
	if string(val) != string(data2) {
		t.Error("Read: expected", string(data2), "got", string(val))
	}

	// Update key2
	s.Write("bucket1", "key2", data1)
	val, _ = s.Read("bucket1", "key2")
	if string(val) != string(data1) {
		t.Error("Update: expected", string(data1), "got", string(val))
	}
}

func testWalk(t *testing.T) {
	// Open a database
	store, _ := NewStore("test.db")
	defer store.Close()
	defer os.Remove("test.db")

	store.CreateBucket("bucket")

	for i := 0; i < 100; i++ {
		store.Write("bucket", fmt.Sprintf("%d", i), nil)
	}

	var buckets []string
	store.Walk(func(key string, val []byte) {
		buckets = append(buckets, key)
	})

	if len(buckets) != 1 || buckets[0] != "bucket" {
		t.Error("Walk: expected one bucket named bucket got", strings.Join(buckets, " "))
	}

	var keys []string
	store.WalkBucket("bucket", func (key string, val []byte) {
		keys = append(keys, key)
	})

	if len(keys) != 100 {
		t.Error("WalkBucket: expected 100 keys got", len(keys))
	}

	var tens []string
	store.WalkPrefix("bucket", "1", func(key string, val []byte) {
		fmt.Println(key)
		tens = append(tens, key)
	})

	if len(tens) != 11 {
		t.Error("WalkPrefix: expected 11 keys got", len(tens))
	}
}

func stringSliceEqual(s1, s2 []string) bool {
	if len(s1) != len(s2) {
		return false
	}

	for i := range s1 {
		if s1[i] != s2[i] {
			return false
		}
	}

	return true
}
