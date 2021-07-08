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
	dbPath   = flag.String("dbPath", "/tmp/fizzbuzz", "path of the database folder")
)

func main() {
	flag.Parse()

	service, err := fizzbuzz.New(*dbPath)
	if err != nil {
		log.Fatalf("cant create service %s", err)
	}

	http.HandleFunc("/", homePage)
	http.HandleFunc("/fizzbuzz", service.FizzBuzzEndpoint)
	http.HandleFunc("/stat", service.Statistics)
	http.HandleFunc("/delete", service.DeleteAll)

	log.Printf("Listening HTTP requests on port: %d", *httpPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *httpPort), nil))
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}
