package store

import (
    "os"
    "testing"
)

var data1 = []byte("Store this.")
var data2 = []byte("Store that.")

func TestStore(t *testing.T) {
    // Open a database in a path that does not exist.
    _, err := NewStore("bad/path/test.db")
    if err == nil {
        t.Error("Expected error while creating new store, got nil")
    }

    // Open a database
    s, err := NewStore("test.db")
    if err != nil {
        t.Error("Unexpected error creating new store", err)
    }
    defer s.Close()
    defer os.Remove("test.db")

    err = s.Create("bucket1", "key1", data)
    if err == nil {
        t.Error("Expected error creating key in non-existent bucket, got nil")
    }

    err = s.CreateBucket("bucket1")
    if err != nil {
        t.Error("Unexpected error creating bucket", err)
    }

    err = s.Create("bucket1", "key1", data1)
    if err != nil {
        t.Error("Unexpected error when creating bucket", err)
    }

    err, val := s.Read("bucket1", "key1")
    if string(val) != string(data1) {
        t.Error("Expected", string(data1), "got", string(val))
    }


}

func TestBucket(t *testing.T) {
    s, err := NewStore("test.db")
    if err != nil {
        t.Error("Unexpected error creating new store", err)
    }
    defer s.Close()
    defer os.Remove("test.db")

    err = s.CreateBucket("bucket")
    if err != nil {
        t.Error("Unexpected error", err)
    }

    err = s.DeleteBucket("bucket")
    if err != nil {
        t.Error("Unexpected error", err)
    }

    err = s.DeleteBucket("nonexistent")
    if err == nil {
        t.Error("Expected an error while deleting bucket, got nil")
    }
}
