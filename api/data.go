package api

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/url"
	"strconv"
	"time"
)

type Data struct {
	UniqueProviders map[string]bool
	PriceListID     string
	ValidUntil      time.Time
	FlightPlans     []Flights
	FormInfo        FormInfo
	Token           string
}

type Flights struct {
	Route         []RouteInfo //Route, that goes with the named flightplan
	Flights       []Provider  //List of flights needed to complete the route
	TravelTime    time.Duration
	Destination   string
	TotalDistance int
	TotalPrice    float64
}

type FormInfo struct {
	MaxPrice          int
	ChosenPrice       int
	MaxDistance       int
	ChosenDistance    int
	MaxTime           int
	ChosenTime        int
	ChosenSort        string
	DepartureLocation string
	Destination       string
}

var allRoutes map[string][]string

//Saves inital data from parsed api-data
func GetData() (Data, error) {
	var data Data

	err := decodeJSON()

	if err != nil {
		return data, err
	}

	allRoutes = mapConnections()

	data.Token = makeToken()
	data.PriceListID = rawData.ID
	data.ValidUntil = rawData.ValidUntil
	data.UniqueProviders = mapProviders()

	return data, err
}

//Makes an unique token
func makeToken() string {
	crutime := time.Now().Unix()
	h := md5.New()
	io.WriteString(h, strconv.FormatInt(crutime, 10))
	token := fmt.Sprintf("%x", h.Sum(nil))
	return token
}

//Maps all unique providers from the data from api
func mapProviders() map[string]bool {

	uniqueProviders := make(map[string]bool)
	for _, leg := range rawData.Legs {
		for _, v := range leg.Providers {
			uniqueProviders[v.Company.Name] = true
		}
	}

	return uniqueProviders
}

//Maps all possible connections between the planets using the data from json
func mapConnections() map[string][]string {

	possibleRoutes := make(map[string][]string)

	for _, v := range rawData.Legs {
		possibleRoutes[v.RouteInfo.From.Name] = append(possibleRoutes[v.RouteInfo.From.Name], v.RouteInfo.To.Name)
	}

	return possibleRoutes
}

//Calculates all the flightplans for all of the possible routes
func makeFlightPlan(routes [][]RouteInfo) ([]Flights, error) {

	var toTheUser []Flights
	var err error
	timestamp := time.Now().UTC()

	for _, route := range routes {

		trips := findTrips(timestamp, route, []Provider{}, [][]Provider{})

		if len(trips) > 0 {

			for _, flightPlan := range trips {
				totalPrice, err := calcTotalPrice(flightPlan)
				if err != nil {
					return toTheUser, err
				}

				tempTrip := Flights{
					route,
					flightPlan,
					calcDuration(flightPlan),
					route[len(route)-1].To.Name,
					calcTotalDistance(route),
					totalPrice,
				}
				toTheUser = append(toTheUser, tempTrip)
			}

		}

	}

	return toTheUser, err
}

//Recursively puts together all valid flightplans for given route
func findTrips(timestamp time.Time, route []RouteInfo, trip []Provider, trips [][]Provider) [][]Provider {
	allTrips := trips
	endTime := timestamp

	for _, v := range rawData.Legs {
		if v.RouteInfo == route[0] {
			for _, validTrip := range v.Providers {
				if endTime.Before(validTrip.FlightStart) {
					tempRoute := append(trip, validTrip)
					if len(route) > 1 {
						allTrips = findTrips(validTrip.FlightEnd, route[1:], tempRoute, allTrips)
					} else {
						var temp []Provider
						for _, v := range tempRoute {
							temp = append(temp, v)
						}
						allTrips = append(allTrips, temp)
					}
				}
			}
			break
		}
	}

	return allTrips
}

//Checks recursivly for all the routes to final destination
func findRoutes(start string, finalDestination string, list [][]string, route []string) [][]string {
	routes := list

	for _, v := range allRoutes[start] { //checks if all connecting routes take to the end of the
		tempRoute := append(route, start)
		if v == finalDestination {
			tempRoute = append(tempRoute, v)
			routes = append(routes, tempRoute)
		} else {
			check := true
			for _, point := range tempRoute { //checks for duplicate locations
				if point == v {
					check = false
				}
			}
			if check { //if no duplicates are found, continues the search for destination
				routes = findRoutes(v, finalDestination, routes, tempRoute)
			}
		}
	}

	return routes
}


//Checks for routes and saves info as Routeinfo
func checkRoutes(form url.Values) [][]RouteInfo {
	var routesAsRouteInfo [][]RouteInfo
	if len(form["from"]) != 0 && len(form["to"]) != 0 {

		valiRoutes := findRoutes(form.Get("from"), form.Get("to"), [][]string{}, []string{})

		//Parses routes from name-based to Routeinfo
		for i := 0; i < len(valiRoutes); i++ {
			var tempRoute []RouteInfo
			for j := 0; j < len(valiRoutes[i])-1; j++ {
				for _, leg := range rawData.Legs {
					if valiRoutes[i][j] == leg.RouteInfo.From.Name && valiRoutes[i][j+1] == leg.RouteInfo.To.Name {
						tempRoute = append(tempRoute, leg.RouteInfo)
						break
					}
					
				}
			}
			routesAsRouteInfo = append(routesAsRouteInfo, tempRoute)
		}
	}

	return routesAsRouteInfo
}
