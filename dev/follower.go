package main

import (
	// "html/template"
	// "io/ioutil"
	"log"
	"net/http"
	// "regexp"
	"fmt"
)

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
	log.Println("[get] key: ", key)
	
	/* now get from persistent storage */

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