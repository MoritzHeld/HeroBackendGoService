package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/handlers"
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
		var maxId int
		for i := 0; i < len(heroes); i++ {
			idint, err := strconv.Atoi(heroes[i].ID)
			if err == nil && idint > 0 && idint > maxId {
				maxId = idint
			}
		}
		newHero.ID = strconv.Itoa(maxId + 1)
	}

	idint, err := strconv.Atoi(newHero.ID)
	if idint < 0 {
		fmt.Fprintf(w, "Hero ID is invalid, below zero. ")
		return
	}

	if newHero.Name == "" {
		newHero.Name = "Name"
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
	//enableCors(&w)
	heroID := mux.Vars(r)["id"]
	var updatedHero hero

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Type in new Data")
	}
	json.Unmarshal(reqBody, &updatedHero)

	for i := range heroes {
		if heroes[i].ID == heroID {
			heroes[i].Name = updatedHero.Name
			heroes[i].Description = updatedHero.Description
			json.NewEncoder(w).Encode(heroes[i])
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
	router.HandleFunc("/heroes", createHero).Methods("PUT")
	router.HandleFunc("/heroes", getAllHeroes).Methods("GET")
	router.HandleFunc("/heroes/{id}", getOneHero).Methods("GET")
	router.HandleFunc("/heroes/{id}", updateHero).Methods("POST")
	router.HandleFunc("/heroes/{id}", deleteHero).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedHeaders([]string{"X-Requested-With"}),
		handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "PATCH", "OPTIONS"}))(router)))

}
