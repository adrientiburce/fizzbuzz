package main

import (
	"fizzbuzz/fizzbuzz"
	"fmt"
	"log"
	"net/http"

	"github.com/namsral/flag"
)

var (
	httpPort = flag.Int("httpPort", 1082, "http port to listen on")
)

func main() {
	flag.Parse()

	service := fizzbuzz.New()

	http.HandleFunc("/", homePage)
	http.HandleFunc("/fizzbuzz", service.FizzBuzzEndpoint)
	http.HandleFunc("/statistics", service.StatisticsEndpoint)

	log.Printf("Listening HTTP requests on port: %d", *httpPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *httpPort), nil))
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}
