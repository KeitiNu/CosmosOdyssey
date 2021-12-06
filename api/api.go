package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

var rawData ApiInfo //global variable for all the data from api
var counter int     //TEST for data requests

type ApiInfo struct {
	ID         string    `json:"id"`
	ValidUntil time.Time `json:"validUntil"`
	Legs       []Legs    `json:"legs"`
}

type Legs struct {
	ID        string     `json:"id"`
	RouteInfo RouteInfo  `json:"routeInfo"`
	Providers []Provider `json:"providers"`
}

type RouteInfo struct {
	ID   string `json:"id"`
	From struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"from"`
	To struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"to"`
	Distance int `json:"distance"`
}

type Provider struct {
	ID      string `json:"id"`
	Company struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"company"`
	Price       float64   `json:"price"`
	FlightStart time.Time `json:"flightStart"`
	FlightEnd   time.Time `json:"flightEnd"`
}

//Parses data from the api
func decodeJSON() error {
	var err error
	apiURL := "https://cosmos-odyssey.azurewebsites.net/api/v1.0/TravelPrices"
	jsonData, err := parseToBytes(apiURL)

	if err != nil {
		return err
	}

	err = json.Unmarshal(jsonData, &rawData)

	if err != nil {
		return err
	}

	//If api has new info, sends that info to the database
	err = apiInfoToDB()
	return err
}

//Parses json data from a given link to []byte
func parseToBytes(link string) ([]byte, error) {
	var jsonData []byte
	response, err := http.Get(link)

	if err != nil {
		log.Printf("createBody: %s", err)
		return jsonData, err
	}

	jsonData, err = io.ReadAll(response.Body)

	if err != nil {
		log.Printf("createBody: %s", err)
		return jsonData, err
	}

	return jsonData, err
}
