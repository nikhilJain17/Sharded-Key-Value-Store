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
	// "log"
	// "net/http"
	// // "regexp"
	"fmt"
	// // "encoding/gob"
	// // "bytes"
	// "github.com/boltdb/bolt"
	// "time" 
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
	followerInit(&follower)
	
}



