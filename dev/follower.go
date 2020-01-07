package main

import (
	// "html/template"
	// "io/ioutil"
	// "log"
	"net/http"
	// "regexp"
	"fmt"
)

func testHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello world")
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		// log.Fatal(err)
	}
	fmt.Println("[get]: ") // log to server console
	fmt.Printf("%+v\n", r.Form)
	for key, values := range r.Form {
		for _, value := range values {
			fmt.Println(key, value)
		}
	}
	// fmt.Println("===========") 	
	// fmt.Println(r) 
	// fmt.Fprintf(w, "get: ")

}

func putHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "put: ")
}

func main() {
	// hello world handler
    http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Welcome to my website!")
	})

	http.HandleFunc("/test", testHandler)
	http.HandleFunc("/get", getHandler)
	http.HandleFunc("/put", putHandler)
	http.ListenAndServe(":8080", nil)
}