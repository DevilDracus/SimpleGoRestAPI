package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

type Log struct {
	DateTime      time.Time `json:"Date"`
	EventType     string    `json:"EventType"`
	Source        string    `json:"Source"`
	VersionNumber string    `json:"VersionNumber"`
	Message       string    `json:"Message"`
	StackTrace    string    `json:"StackTrace"`
	CalledMethod  string    `json:"CalledMethod"`
}

type Logs []Log

func postLog(w http.ResponseWriter, r *http.Request) {
	var log Log
	json.NewDecoder(r.Body).Decode(&log)
	fmt.Fprintf(w, "Post Endpoint Hit")
	writeJSON(log)
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/log", postLog).Methods("POST")
	log.Fatal(http.ListenAndServe(":8081", myRouter))
}

func writeJSON(data Log) {
	var logs []Log

	filename := "logs.json"
	err := checkFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	file, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal(file, &logs)

	logs = append(logs, data)

	dataBytes, err := json.MarshalIndent(logs, "", " ")
	if err != nil {
		log.Println("Unable to create JSON file")
		return
	}

	_ = ioutil.WriteFile("logs.json", dataBytes, 0644)
}

func checkFile(filename string) error {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		_, err := os.Create(filename)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	handleRequests()
}
