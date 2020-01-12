/* The code for the leader in RAFT. 

> the leader is constantly sending heartbeats to followers?

Once the leader is elected...
> the client sends a change to the leader
> change is appended to leader's log
> change is sent to followers on next heartbeat
	> client puts it on their log
> response is sent to client
	> an entry is "committed" if a majority of followers ack that shite
	> else it's an abort 

> when theres a new leader elected, we roll back uncommitted logs

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

type LogEntry struct {
	UID string
	Timestamp time.Time
	Status string // abort or committed
	Req Request// @todo
}

func setupFollowers(numFollowers int) {
	var urls [3]string
	urls[0] = ":8080"
	urls[1] = ":8081"
	urls[2] = ":8082"
	var db [3]string
	db[0] = "kvstore0.db"
	db[1] = "kvstore1.db"
	db[2] = "kvstore2.db"
	for i := 0; i < 3; i++ {
		follower := Follower {
			UID : string(i),
			DBFilename : db[i],
			URL : urls[i],
		}	
		followerInit(&follower) // start up the follower server to listen to http requests
	}
}

func main() {
	// get follower servers up and running
	fmt.Println("hello")
	go setupFollowers(1)
	// var urls [3]string
	// urls[0] = ":8080"
	// urls[1] = ":8081"
	// urls[2] = ":8082"
	// var db [3]string
	// db[0] = "kvstore0.db"
	// db[1] = "kvstore1.db"
	// db[2] = "kvstore2.db"
	// for i := 0; i < 3; i++ {
	// 	follower := Follower {
	// 		UID : string(i),
	// 		DBFilename : db[i],
	// 		URL : urls[i],
	// 	}	
	// 	followerInit(&follower) // start up the follower server to listen to http requests
	// }

	// set up log for leader
	log := make ([]LogEntry, 1)
	log = append(log, LogEntry{UID : "leader", Timestamp : time.Now(), Status : "initial", })
	
	// > set up server stuff for leader
	heartbeatChannel := make(chan Request)

	// server loop
	setupLeaderServer(heartbeatChannel)
	go http.ListenAndServe(":5000", nil)

	// wait until follower is up to send heartbeats
	time.Sleep(5 * time.Second)
	// heartbeat loop 
	for true {
		heartbeat(heartbeatChannel)
	}
	
}


func setupLeaderServer(heartbeatChannel chan Request) {

	// Once the leader is elected...
	// > the client sends a change to the leader
	// > change is appended to leader's log
	// > change is sent to followers on next heartbeat
	// > response is sent to client
	// 	> an entry is "committed" if a majority of followers ack that shite
	// 	> else it's an abort 
		

	// client http handlers
	http.HandleFunc("/getLeader", func (w http.ResponseWriter, r *http.Request) {
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
	
	// send to all followers
	var urls [3]string
	urls[0] = "http://127.0.0.1:8080/heartbeat"
	urls[1] = "http://127.0.0.1:8081/heartbeat"
	urls[2] = "http://127.0.0.1:8082/heartbeat"
	for i := 0; i < 3; i++ {
		resp, err := http.PostForm(urls[i], form)
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
	}

			
	// sleep
	time.Sleep(5 * time.Second)
	
}

// follower:
// curl -d "value=me&key=name" -H "Content-Type: application/x-www-form-urlencoded" -X POST http://127.0.0.1:8080/put

// leader:
// curl http://127.0.0.1:8001/test