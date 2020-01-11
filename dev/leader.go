/* The code for the leader in RAFT. 

> the leader is constantly sending heartbeats to followers?

Once the leader is elected...
> the client sends a change to the leader
> change is appended to leader's log
> change is sent to followers on next heartbeat
> response is sent to client
	> an entry is "committed" if a majority of followers ack that shite
	> else it's an abort or some shite

todo: something for network partitions and elections or whatever, i think just implement elections and things will be fine
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
	// "strings"
	"net/url"
)

type Follower struct {
	UID string
	DBFilename string
	URL string
}

type Request struct {
	Type string // GET, PUT, DELETE (@todo make this an enum?)
	Kvpair KVPair // for put 
	Key string // for get, delete
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

	heartbeatChannel := make(chan Request)

	// server loop
	setupLeaderServer(heartbeatChannel)
	go http.ListenAndServe(":5000", nil)

	// wait until follower is up to send heartbeats
	time.Sleep(2 * time.Second)
	// heartbeat loop 
	for true {
		heartbeat(heartbeatChannel)
	}
	
}


func setupLeaderServer(heartbeatChannel chan Request) {
	// client http handlers
	http.HandleFunc("/getLeader", func (w http.ResponseWriter, r *http.Request) {
		// @todo
		// log.Println("get on Leader")
		fmt.Fprintf(w, "getting")

		// send whatever client sent thru the heartbeatChannel
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
		log.Println("[leaderGet] key:", key)

		heartbeatChannel <- Request{Type : "GET", Kvpair : KVPair{Key : key}, Key : key,} // @todo sending empty kvpair?
	})


	http.HandleFunc("/putLeader", func (w http.ResponseWriter, r *http.Request) {
		// @todo
		fmt.Fprintf(w, "putting")

		// send whatever client sent thru the heartbeatChannel
		err := r.ParseForm()
		if err != nil {
			log.Fatal(err)
		}	
		var key string 
		var value string 
		for formKey, formValues := range r.Form {
			for _, formVal := range formValues { // @todo make sure length of values is one
				if string(formKey) == "key" {
					key = formVal
				}
				if string(formKey) == "value" {
					value = formVal
				}
			}
		}
		log.Println("[leaderPut] key:", key, "value:", value)

		heartbeatChannel <- Request{Type : "PUT", Kvpair : KVPair{Key : key, Value : value}, Key : key,} // @todo sending empty kvpair?
	})
}

func heartbeat(heartbeatChannel chan Request) {	
	// hc := http.Client{}
	form := url.Values{}

	// appendArgs := false

	// non blocking channels for the heartbeat!
	select {
	case leaderReq := <- heartbeatChannel:
		log.Println("made it into the channel:", leaderReq)
		form.Add("key", leaderReq.Kvpair.Key)
		form.Add("value", leaderReq.Kvpair.Value)
		form.Add("type", leaderReq.Type)
		
	default:
		log.Println("nothing yet")
	}

	log.Println("heartbeat form:", form)

	// send post request to heartbeat WITH request data if applicable @todo
	form.Add("heartbeat", "true")
	
	resp, err := http.PostForm("http://127.0.0.1:8080/heartbeat", form)
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
		log.Println("[heartbeat reply in leader]:", bodyString)
	}
			
	// sleep
	time.Sleep(5 * time.Second)
	
}

// follower:
// curl -d "value=me&key=name" -H "Content-Type: application/x-www-form-urlencoded" -X POST http://127.0.0.1:8080/put

// leader:
// curl http://127.0.0.1:8001/test