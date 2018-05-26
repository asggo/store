package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/asggo/store"
)

const (
	server         = "127.0.0.1"
	port           = "5000"
	buckTmplString = `
<html>
<head></head>
<body>
<h1>Buckets</h1>
<u>
{{range .Buckets}}
<li><a href="/get/{{.}}/">{{.}}</a></li>
{{end}}
</u>
</body>
</html>`

	keyTmplString = `
<html>
<head></head>
<body>
<h1>{{.Bucket}} - Keys</h1>
<u>
{{range .Keys}}
<li><a href="/get/{{$.Bucket}}/{{.}}/">{{.}}</a></li>
{{end}}
</u>
</body>
</html>`

	valTmplString = `
<html>
<head><title>Store Web View</title></head>
<body>
<h1>{{.Bucket}} - {{.Key}}</h1>
<pre>{{.Value}}</pre>
</body>
</html>`
)

var (
	dbfile   = ""
	buckTmpl = template.Must(template.New("bucket").Parse(buckTmplString))
	keyTmpl  = template.Must(template.New("key").Parse(keyTmplString))
	valTmpl  = template.Must(template.New("val").Parse(valTmplString))
)

type Buckets struct {
	Buckets []string
}

type Keys struct {
	Bucket string
	Keys   []string
}

type Value struct {
	Bucket string
	Key    string
	Value  string
}

func conn() *store.Store {
	db, err := store.NewStore(dbfile)
	if err != nil {
		log.Printf("Could not open connection to %s: %s\n", dbfile, err)
		return nil
	}

	return db
}

func getBuckets() []string {
	db := conn()

	if db == nil {
		return nil
	}

	defer db.Close()

	items, err := db.AllBuckets()
	if err != nil {
		log.Printf("Could not retrieve bucket list: %s\n", err)
		return nil
	}

	return items
}

func getKeys(bucket string) []string {
	db := conn()

	if db == nil {
		return nil
	}

	defer db.Close()

	items, err := db.AllKeys(bucket)
	if err != nil {
		log.Printf("Could not retrieve keys from bucket %s: %s\n", bucket, err)
		return nil
	}

	return items
}

func getValue(bucket, key string) string {
	db := conn()

	if db == nil {
		return ""
	}

	defer db.Close()

	val := db.Read(bucket, key)
	return string(val)
}

func findBuckets(query string) []string {
	db := conn()

	if db == nil {
		return nil
	}

	defer db.Close()

	items, err := db.FindBuckets(query)
	if err != nil {
		log.Printf("Could not find buckets matching %s: %s\n", query, err)
		return nil
	}

	return items
}

func findKeys(bucket, query string) []string {
	db := conn()

	if db == nil {
		return nil
	}

	defer db.Close()

	items, err := db.FindKeys(bucket, query)
	if err != nil {
		log.Printf("Could not find keys matching %s in bucket %s: %s\n", bucket, query, err)
		return nil
	}

	return items
}

func get(w http.ResponseWriter, r *http.Request) {
	args := strings.Split(r.URL.Path, "/")

	switch len(args) {
	case 3:
		buckets := getBuckets()
		buckTmpl.Execute(w, &Buckets{Buckets: buckets})
	case 4:
		keys := getKeys(args[2])
		keyTmpl.Execute(w, &Keys{Bucket: args[2], Keys: keys})
	case 5:
		val := getValue(args[2], args[3])
		valTmpl.Execute(w, &Value{Bucket: args[2], Key: args[3], Value: val})
	default:

	}
}

func find(w http.ResponseWriter, r *http.Request) {
	args := strings.Split(r.URL.Path, "/")
fmt.Println(args, len(args))
	switch len(args) {
	case 2:
		buckets := getBuckets()
		buckTmpl.Execute(w, &Buckets{Buckets: buckets})
	case 3:
		buckets := findBuckets(args[2])
		buckTmpl.Execute(w, &Buckets{Buckets: buckets})
	case 4:
		keys := findKeys(args[2], args[3])
		keyTmpl.Execute(w, &Keys{Bucket: args[2], Keys: keys})
	default:
	}
}

// Setup our HTTP server and route handlers.
func main() {

	if len(os.Args) < 2 {
		fmt.Println("Usage: kvweb.go dbfile")
	}

	dbfile = os.Args[1]

	http.HandleFunc("/get/", get)
	http.HandleFunc("/find/", find)

	err := http.ListenAndServe(fmt.Sprintf("%s:%s", server, port), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
