package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type ScheduleTag struct {
	Id   int    `json:"id"`
	Name string `json:"schedule"`
}

type PlaybackLog struct {
	Id          int         `json:"id"`
	ScheduleTag ScheduleTag `json:"scheduleTag"`
	Screen      string      `json:"schedule"`
	Cpl         string      `json:"orderId"`
}

var Queue []PlaybackLog

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "AgentQL homepage endpoint hit")
}

func receiveLogs(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: receiveLogs")

	responseData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(responseData))

	// pushes the received logs into queue
	Ack := make([]PlaybackLog, 10)

	json.Unmarshal(responseData, &Ack)

	Queue = append(Queue, Ack...)
}

func handleRequests() {
	http.HandleFunc("/", homePage)

	// add receiveLogs route and map it to our receiveLogs function like so
	http.HandleFunc("/receiveLogs", receiveLogs)
	log.Fatal(http.ListenAndServe(":8083", nil))
}

func main() {
	// initialise agentql apis
	go handleRequests()

	// end-less process
	for {
		rand.Seed(time.Now().UnixNano())
		time.Sleep(time.Second * time.Duration(rand.Intn(10)))

		// sends acknowledgement of log delivery
		if len(Queue) != 0 {
			jsonValue, _ := json.Marshal(Queue[0].Id)

			_, err := http.Post("http://localhost:8082/ackFromAgentQL", "application/json", bytes.NewBuffer(jsonValue))

			if err != nil {
				fmt.Print(err.Error())
			}
		}

		if len(Queue) > 1 {
			Queue = Queue[1:]
		} else {
			Queue = make([]PlaybackLog, 0)
		}
	}
}
