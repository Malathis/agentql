package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func saveLogs(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: saveLogs")

	responseData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(responseData))
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "AgentQL homepage endpoint hit")
}

func handleRequests() {
	http.HandleFunc("/", homePage)

	// add our logs route and map it to our
	// saveLogs function like so
	http.HandleFunc("/saveLogs", saveLogs)
	log.Fatal(http.ListenAndServe(":8083", nil))
}

func main() {
	handleRequests()
}
