package main

import (
	"encoding/json"
	"fizzbuzz/fizzbuzz"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/namsral/flag"
)

var (
	httpPort = flag.Int("httpPort", 1082, "http port to listen on")
)

func main() {
	flag.Parse()

	http.HandleFunc("/", homePage)
	http.HandleFunc("/fizzbuzz", fizzBuzz)
	log.Printf("Listening HTTP requests on port: %d", *httpPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *httpPort), nil))
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func fizzBuzz(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)

	var fizzBuzz fizzbuzz.FizzBuzz
	err := json.Unmarshal(reqBody, &fizzBuzz)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to decode request (%s)", err), http.StatusInternalServerError)
		return
	}
	if fizzBuzz.Limit == 0 {
		http.Error(w, "fizzbuzz limit empty", http.StatusBadRequest)
	} else if fizzBuzz.Int1 == 0 || fizzBuzz.Int2 == 0 {
		http.Error(w, "fizzbuzz int1 or int2 empty", http.StatusBadRequest)
	} else if fizzBuzz.Str1 == "" || fizzBuzz.Str2 == "" {
		http.Error(w, "fizzbuzz str1 or str2 empty", http.StatusBadRequest)
	}

	res := fizzBuzz.FizzBuzz()
	fmt.Fprint(w, res)
}
