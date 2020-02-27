package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome home!")
}

type hero struct {
	ID          string `json:"ID"`
	Name        string `json:"Name"`
	Description string `json:"Description"`
}

type allHeroes []hero

var heroes = allHeroes{
	{
		ID:          "1",
		Name:        "Moritz",
		Description: "held1",
	},
	{
		ID:          "2",
		Name:        "Fabian",
		Description: "held2",
	},
	{
		ID:          "3",
		Name:        "Sven",
		Description: "held3",
	},
}

func createHero(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	var newHero hero
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Type in the Data of the new Hero")
	}

	json.Unmarshal(reqBody, &newHero)

	_, parseErr := strconv.Atoi(newHero.ID)
	if parseErr != nil {
		fmt.Fprintf(w, "Hero ID is invalid. ")
		return
	}

	idint, err := strconv.Atoi(newHero.ID)
	if idint < 0 {
		fmt.Fprintf(w, "Hero ID is invalid, below zero. ")
		return
	}

	heroes = append(heroes, newHero)
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(newHero)
}

func getOneHero(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	eventID := mux.Vars(r)["id"]

	for _, singleHero := range heroes {
		if singleHero.ID == eventID {
			json.NewEncoder(w).Encode(singleHero)
		}
	}
}

func getAllHeroes(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	json.NewEncoder(w).Encode(heroes)
}

func updateHero(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	eventID := mux.Vars(r)["id"]
	var updatedHero hero

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Type in new Data")
	}
	json.Unmarshal(reqBody, &updatedHero)

	for i, singleHero := range heroes {
		if singleHero.ID == eventID {
			singleHero.Name = updatedHero.Name
			singleHero.Description = updatedHero.Description
			heroes = append(heroes[:i], singleHero)
			json.NewEncoder(w).Encode(singleHero)
		}
	}
}

func deleteHero(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	HeroID := mux.Vars(r)["id"]

	for i, singleHero := range heroes {
		if singleHero.ID == HeroID {
			heroes = append(heroes[:i], heroes[i+1:]...)
			fmt.Fprintf(w, "The Hero with ID %v has been deleted successfully", HeroID)
		}
	}
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	router.HandleFunc("/heroes", createHero).Methods("POST")
	router.HandleFunc("/heroes", getAllHeroes).Methods("GET")
	router.HandleFunc("/heroes/{id}", getOneHero).Methods("GET")
	router.HandleFunc("/heroes/{id}", updateHero).Methods("PATCH")
	router.HandleFunc("/heroes/{id}", deleteHero).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", router))
}
