package store

import (
	"fmt"
	"os"
	"sort"
	"testing"
)

func TestStore(t *testing.T) {
	testNewStore(t)
	testBucket(t)
	testKey(t)
	testSearch(t)
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
		t.Fatal("Unexpected error", err)
	}

	err = s.DeleteBucket("bucket")
	if err != nil {
		t.Fatal("Unexpected error", err)
	}

	err = s.DeleteBucket("nonexistent")
	if err == nil {
		t.Fatal("Expected an error while deleting bucket, got nil")
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
		t.Fatal("Write: unexpected error", err)
	}

	s.Write("bucket1", "key2", data2)

	// Read keys
	val := s.Read("bucket1", "key1")
	if string(val) != string(data1) {
		t.Error("Read: expected", string(data1), "got", string(val))
	}

	val = s.Read("bucket1", "key2")
	if string(val) != string(data2) {
		t.Error("Read: expected", string(data2), "got", string(val))
	}

	// Update key2
	s.Write("bucket1", "key2", data1)
	val = s.Read("bucket1", "key2")
	if string(val) != string(data1) {
		t.Error("Update: expected", string(data1), "got", string(val))
	}
}

func testSearch(t *testing.T) {
	// Open a database
	store, _ := NewStore("test.db")
	defer store.Close()
	defer os.Remove("test.db")

	store.CreateBucket("bucket")

	var all []string
	var even []string

	for i := 0; i < 100; i++ {
		s := fmt.Sprintf("odd%d", i)

		if i%2 == 0 {
			s = fmt.Sprintf("even%d", i)
			even = append(even, s)
		}

		all = append(all, s)

		store.Write("bucket", s, nil)
	}

	sort.Strings(all)
	sort.Strings(even)

	// Get all keys
	allKeys, err := store.AllKeys("bucket")
	if err != nil {
		t.Error("All Keys: unexpected error", err)
	}
	sort.Strings(allKeys)

	if !stringSliceEqual(allKeys, all) {
		t.Error("Find Keys: expected", all, "got", allKeys)
	}

	evenKeys, err := store.FindKeys("bucket", "even")
	if err != nil {
		t.Error("Find Keys: unexpected error", err)
	}
	sort.Strings(evenKeys)

	if !stringSliceEqual(evenKeys, even) {
		t.Error("Find Keys: expected", even, "got", evenKeys)
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
