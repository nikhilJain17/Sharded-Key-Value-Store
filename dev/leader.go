/* The code for the leader in RAFT. 

> the leader is constantly sending heartbeats to followers?

Once the leader is elected...
> the client sends a change to the leader
> change is appended to leader's log
> change is sent to followers on next heartbeat
> response is sent to client
	> an entry is "committed" if a majority of followers ack that shite
	> else it's an abort or some shite

todo: some shit for network partitions and elections or whatever, i think just implement elections and things will be fine
*/

package main

import (
	// // "html/template"
	// // "io/ioutil"
	"log"
	"net/http"
	// // "regexp"
	"fmt"
	// // "encoding/gob"
	// // "bytes"
	// "github.com/boltdb/bolt"
	"time" 
	"strings"
	"net/url"
)

type Follower struct {
	UID string
	DBFilename string
	URL string
}

func main() {
	fmt.Println("hello")
	// var followers []Follower
	follower := Follower {
		UID : "a humble test",
		DBFilename : "kvstore.db",
		URL : ":8080",
	}
	fmt.Println(follower)
	go followerInit(&follower) // start up the follower server to listen to http requests
	// heartbeat loop 
	fmt.Println("hi")

	hc := http.Client{}
	for true {
		// send post request to heartbeat

		form := url.Values{}
		form.Add("heartbeat", "true")
		req, err := http.NewRequest("POST", "http://127.0.0.1:8080/heartbeat", strings.NewReader(form.Encode()))
		if err != nil {
			log.Fatal(err)
		}
		resp, err := hc.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(resp) // print the heartbeat ack

		
		// sleep
		time.Sleep(5 * time.Second)
	} 
	
}

// hc := http.Client{}
// req, err := http.NewRequest("POST", APIURL, nil)

// form := url.Values{}
// form.Add("heartbeat", "true")
// req, err := http.NewRequest("POST", url, strings.NewReader(form.Encode()))

// req.PostForm = form

// glog.Info("form was %v", form)
// resp, err := hc.Do(req)

// req, err := http.NewRequest("POST", url, strings.NewReader(form.Encode()))

// curl -d "heartbeat=true" -H "Content-Type: application/x-www-form-urlencoded" -X POST http://127.0.0.1:8080/heartbeat



