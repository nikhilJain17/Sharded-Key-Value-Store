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
		
		select {
		case resp :=  <- heartbeatChannel : // what is the result of the heartbeat
			if resp.Status == "ABORT" {
				fmt.Fprintf(w, "ABORTED: " + string(resp.Type) + " " + string(resp.Kvpair.Key) + ", " + string(resp.Kvpair.Value))
			} else if resp.Status == "COMMIT" {
				fmt.Fprintf(w, "COMMITTED: " + string(resp.Type) + " " + string(resp.Kvpair.Key) + ", " + string(resp.Kvpair.Value))
			}
		default:
			// non blocking op
			log.Println("waiting")
		}
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

		// send request
		heartbeatChannel <- Request{Type : "PUT", Kvpair : KVPair{Key : key, Value : value}, Key : key,} // @todo sending empty kvpair?
	
		// wait for response		
		resp :=  <- heartbeatChannel  // what is the result of the heartbeat
		if resp.Status == "ABORT" {
			fmt.Fprintf(w, "ABORTED: " + string(resp.Type) + " " + string(resp.Kvpair.Key) + ", " + string(resp.Kvpair.Value))
		} else if resp.Status == "COMMIT" {
			fmt.Fprintf(w, "COMMITTED: " + string(resp.Type) + " " + string(resp.Kvpair.Key) + ", " + string(resp.Kvpair.Value))
		}
		// select {
		// case resp :=  <- heartbeatChannel : // what is the result of the heartbeat
		// 	log.Println("in leader server: ", resp.Status)
		// 	if resp.Status == "ABORT" {
		// 		fmt.Fprintf(w, "ABORTED: " + string(resp.Type) + " " + string(resp.Kvpair.Key) + ", " + string(resp.Kvpair.Value))
		// 	} else if resp.Status == "COMMIT" {
		// 		fmt.Fprintf(w, "COMMITTED: " + string(resp.Type) + " " + string(resp.Kvpair.Key) + ", " + string(resp.Kvpair.Value))
		// 	}
		// default:
		// 	// non blocking op
		// 	log.Println("waiting")
		// }


	
	})
}

func heartbeat(heartbeatChannel chan Request) {	
	// hc := http.Client{}
	form := url.Values{}

	// appendArgs := false
	var key string
	var value string
	var reqType string

	// non blocking channels for the heartbeat!
	select {
	case leaderReq := <- heartbeatChannel:
		log.Println("made it into the channel:", leaderReq)
		form.Add("key", leaderReq.Kvpair.Key)
		form.Add("value", leaderReq.Kvpair.Value)
		form.Add("type", leaderReq.Type)
		key = leaderReq.Kvpair.Key
		value = leaderReq.Kvpair.Value
		reqType = leaderReq.Type
		
	default:
		log.Println("nothing yet")
	}

	log.Println("[heartbeat leader] sending form:", form)

	// send post request to heartbeat WITH request data if applicable @todo
	form.Add("heartbeat", "true")
	
	// send to all followers
	var urls [3]string
	urls[0] = "http://localhost:8000/heartbeat"
	urls[1] = "http://16a2a699.ngrok.io/heartbeat"
	urls[2] = "http://5db046fd.ngrok.io/heartbeat"

	numAcks := 0

	for i := 0; i < 3; i++ {
		resp, err := http.PostForm(urls[i], form)
		resp.ParseForm()
		log.Println("gucci gang", resp)
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
			numAcks++
		}
	}

	log.Println("numAcks", numAcks)

	if (numAcks > 3 / 2) && (key != "") {
		// majority of followers committed
		// need to send ack to client
		req := Request{
			Kvpair : KVPair {Key: key, Value: value}, 
			Status:"COMMIT",
			Type : reqType,
		} 
		heartbeatChannel <- req
		log.Println(req)
		
	} else if (numAcks < 3 / 2) && (key != ""){
		// not committed -- abort!
		// need to send a delete to followers
		// need to send ABORT to client
		req := Request{
			Kvpair : KVPair {Key: key, Value: value}, 
			Status:"ABORT",
			Type : reqType,
		} 		
		// heartbeatChannel <- req
		log.Println(req)

	}

	// sleep
	time.Sleep(5 * time.Second)
	
}

// follower:
// curl -d "value=me&key=name" -H "Content-Type: application/x-www-form-urlencoded" -X POST http://127.0.0.1:8080/put

// leader:
// curl http://127.0.0.1:8001/test