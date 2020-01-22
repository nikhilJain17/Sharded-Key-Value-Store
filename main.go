package main

import (
	"time"
	"net/http"
	"fmt"
)

type Follower struct {
	UID string
	DBFilename string
	URL string
}

// type Leader struct

type KVPair struct {
	Key string
	Value string
}

type Request struct {
	Type string // GET, PUT, DELETE (@todo make this an enum?)
	Kvpair KVPair // for put 
	Key string // for get, delete
	Status string // COMMIT, ABORT
}

type LogEntry struct {
	UID string
	Timestamp time.Time
	Status string // abort or committed
	Req Request// @todo
}

func setupFollowers(numFollowers int) {
	// @todo use numFollowers
	var urls [3]string
	urls[0] = "localhost:8000"
	urls[1] = "localhost:8001"
	urls[2] = "localhost:8002"
	var db [3]string
	db[0] = "kvstore0.db"
	db[1] = "kvstore1.db"
	db[2] = "kvstore2.db"

	/* @todo 
		1. make a new binary for the follower
		2. make a new binary for the leader 
		3. exec the follower binary here
	*/
	for i := 0; i < 3; i++ {
		follower := Follower {
			UID : string(i),
			DBFilename : db[i],
			URL : urls[i], // + string(uid),
		}	
		go followerInit(&follower) // start up the follower server to listen to http requests
	}
}


func main() {
	// get follower servers up and running
	fmt.Println("hello")
	setupFollowers(1)
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