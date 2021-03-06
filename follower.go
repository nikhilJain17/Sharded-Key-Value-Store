package main

import (
	// "html/template"
	// "io/ioutil"
	"log"
	"net/http"
	// "regexp"
	"fmt"
	// "encoding/gob"
	// "bytes"
	"github.com/boltdb/bolt"
	"time"
)

var db *bolt.DB

func testHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello world")
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	/* do some input validation etc here probably */
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	var key string 
	for formKey, formValues := range r.Form {
		for _, formVal := range formValues { // @todo make sure length of values is one
			if string(formKey) == "key" {
				key = formVal
			}
		}
	}
	log.Println("[followerGet] key: ", key)
	
	/* now get from persistent storage */

	// access db
	// db, err := bolt.Open("kvstore.db", 0600, nil)

	// read
	var val []byte
	db.View(func(transaction *bolt.Tx) error {
		bucket := transaction.Bucket([]byte("kvbucket"))
		val = bucket.Get( []byte(key) ) // @todo copy the bytes over to a new slice...because val is dead when xact is dead
		fmt.Fprintf(w, string(val))
		return nil
	})

}

func putHandler(w http.ResponseWriter, r *http.Request) {
	/* do some input validation etc here probably */
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	var key string 
	var val string
	for formKey, formValues := range r.Form {
		for _, formVal := range formValues { // @todo make sure length of values is one
			if string(formKey) == "key" {
				key = formVal
			}
			if string(formKey) == "value" {
				val = formVal
			}
		}
	}
	log.Println("[put] key:", key, ", value:", val)

	// @todo some validation here also

	/* store in persistent storage probably */

	// @todo encoding to bytes? (look into gob)

	// var buf bytes.Buffer
	// encoder := gob.NewEncoder(&buf)
	// kvpair := KVPair{
	// 	Key: key,
	// 	Value: val,
	// }
	// encoder.Encode(kvpair)
	// encodedBytes := buf.Bytes()

	// keyBytes := 

	// fmt.Println("sankruth:", encodedBytes)


	// store in a db? sqlite, boltdb...
	// the point of this project is to do the distributed part, not really the kv store part.
	// i don't want to implement a database, i want to implement the distributed algos
	// so we'll just roll with boltdb instead of storing to a sorted file or something lol

	db.Update(func(transaction *bolt.Tx) error {
		bucket, err := transaction.CreateBucketIfNotExists([]byte("kvbucket"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		} 
		err = bucket.Put([]byte(key), []byte(val))
		return err
	})

}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	// @todo
	
}

func heartbeatHandler(w http.ResponseWriter, r *http.Request) {
	// @todo 
	fmt.Fprintf(w, "ack\n")
	
	// if get, then do get but if put, then do put
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}
	var key string 
	var value string 
	var reqType string 

	for formKey, formValues := range r.Form {
		for _, formVal := range formValues { // @todo make sure length of values is one
			if string(formKey) == "key" {
				key = formVal
			}
			if string(formKey) == "value" {
				value = formVal
			}
			if string(formKey) == "type" {
				reqType = formVal
			}
		}
	}

	var val string
	if reqType == "GET" {
		val = getHandler(w, r)
		fmt.Fprintf(w, val)
	} else if reqType == "PUT" {
		putHandler(w, r)
	}

	log.Println("[heartbeat client] key:", key, "value:", value, "reqType:", reqType)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello")	
}

func followerInit(f *Follower) {

	/* set up boltdb */
	var err error
	// @todo each client needs it's own db...
	db, err = bolt.Open(f.DBFilename, 0600, &bolt.Options{Timeout: 2 * time.Second})
	if err != nil {log.Fatal("error opening boltdb: ", err, "client: ", f)}
	defer db.Close() // defer is pretty handy

	log.Println("setting up follower:", f)


	server := http.NewServeMux()

	// hello world handler
    server.HandleFunc("/", helloHandler)

	server.HandleFunc("/test", testHandler)
	server.HandleFunc("/get", getHandler)
	server.HandleFunc("/put", putHandler)
	server.HandleFunc("/delete", deleteHandler)
	server.HandleFunc("/heartbeat", heartbeatHandler)

	http.ListenAndServe(f.URL, server)
}

/*
Some random notes:

- forget using rpc and custom (un)marshaling etc, just use http to make life easier
- use boltdb for persistent storage, focus here is to learn about the distributed part, not the kv store part
- in memory log for now, will change later 
*/