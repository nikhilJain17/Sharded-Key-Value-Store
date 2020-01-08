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
	"io/ioutil"
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
	
	// do we want the heartbeat to be a separate goroutine? 
	// the issue is we need a handler for client requests
	// and we can't have the heartbeat block the thread...so yes?
	// i also dont know how the handler thing works, like if i set something as a handler,
	// will it run in a new thread each time?
	// @todo look into that.

	// server loop
	setupLeaderServer()
	go http.ListenAndServe(":5000", nil)

	
	// heartbeat loop 
	for true {
		heartbeat()
	}
	
}


func setupLeaderServer() {
	// client http handlers
	http.HandleFunc("/sankruth", func (w http.ResponseWriter, r *http.Request) {
		// @todo
		log.Println("yo quit calling on my endpoints dawg")
		fmt.Fprintf(w, "hello world")

	})
}

func heartbeat() {
	hc := http.Client{}
	
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
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString := string(bodyBytes)
		log.Println(bodyString)
	}
			
	// sleep
	time.Sleep(5 * time.Second)
	
}

// follower:
// curl -d "value=me&key=name" -H "Content-Type: application/x-www-form-urlencoded" -X POST http://127.0.0.1:8080/put

// leader:
// curl http://127.0.0.1:8001/test
