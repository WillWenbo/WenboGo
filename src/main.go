package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	// elastic search is a backup solutioin
	// "github.com/elastic/go-elasticsearch/v8"
	"github.com/gorilla/mux"
	"gopkg.in/yaml.v3"
)

// func getElasticSearchClient() {
// 	client, err := elasticsearch.NewDefaultClient()
// 	log.Println(elasticsearch.Version)
// 	return client, err
// }

type event struct {
	Title       string       `json:"Title"`
	Version     string       `json:"Version"`
	Maintainers []maintainer `json:"Maintainers"`
	Company     string       `json:"Company"`
	Website     string       `json:"Website"`
	Source      string       `json:"Source"`
	License     string       `json:"License"`
	Description string       `json:"Description"`
}

type eventSearchParam struct {
	Title            string
	Version          string
	MaintainersEmail string
	maintainersName  string
	Company          string
	Website          string
	Source           string
	License          string
	Description      string
}

type maintainer struct {
	Name  string
	Email string
}

type allEvents []event

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/event", createEvent).Methods("POST")
	router.HandleFunc("/events", getEventsByParams).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func createEvent(w http.ResponseWriter, r *http.Request) {

	var newEvent event
	// validator
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter request data")
	}

	yaml.Unmarshal(reqBody, &newEvent)

	if validateReq(newEvent) != "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(validateReq(newEvent))

	} else {

		if !isEventExist(newEvent) {
			SaveEvent(newEvent)
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newEvent)
	}
}

// Data dedup, check if event exists in DB.
func isEventExist(newEvent event) bool {
	var resultList = searchEventByField(newEvent)
	return len(resultList) != 0
}

func getEventsByParams(w http.ResponseWriter, r *http.Request) {

	var eventParams eventSearchParam

	var url = r.URL
	print(url)

	if r.URL.Query().Get("Title") != "" {
		eventParams.Title = r.URL.Query().Get("Title")
	}

	if r.URL.Query().Get("Version") != "" {
		eventParams.Version = r.URL.Query().Get("Version")
	}

	if r.URL.Query().Get("MaintainersEmail") != "" {
		eventParams.MaintainersEmail = r.URL.Query().Get("MaintainersEamil")
	}

	if r.URL.Query().Get("MaintainersName") != "" {
		eventParams.maintainersName = r.URL.Query().Get("MaintainersName")
	}

	if r.URL.Query().Get("Company") != "" {
		eventParams.Company = r.URL.Query().Get("Company")
	}

	if r.URL.Query().Get("Website") != "" {
		eventParams.Website = r.URL.Query().Get("Website")
	}

	if r.URL.Query().Get("Source") != "" {
		eventParams.Source = r.URL.Query().Get("Source")
	}

	if r.URL.Query().Get("License") != "" {
		eventParams.License = r.URL.Query().Get("License")
	}

	if r.URL.Query().Get("Description") != "" {
		eventParams.Description = r.URL.Query().Get("Description")
	}

	var eventIds = searchEventByField(eventParams)
	eventIds = append(eventIds, 0)
}
