package api

import (
	"fmt"
	"log"
	"strconv"
	"time"
)

//Finds the value of the most expensive route
func findMostExpensive(trips []Flights) int {
	var max float64

	for i, flight := range trips {
		if i == 0 || flight.TotalPrice > max {
			max = flight.TotalPrice
		}
	}
	return int(max + 1)
}

//Finds the value of the longest route
func findLongest(trips []Flights) int {
	var max int

	for i, trip := range trips {
		if i == 0 || trip.TotalDistance > max {
			max = trip.TotalDistance
		}
	}
	return max
}

//Calculates the Duration of the trip
func calcDuration(trip []Provider) time.Duration {
	
	// return trip[len(trip)-1].FlightEnd.Sub(trip[0].FlightStart)
	t:= trip[len(trip)-1].FlightEnd.Sub(trip[0].FlightStart)
	return t
}


//Calculates total price of the flightplan
func calcTotalPrice(data []Provider) (float64, error) {
	var totalPrice float64
	var err error
	for _, flight := range data {
		totalPrice += flight.Price
		i := fmt.Sprintf("%.2f", totalPrice)
		totalPrice, err = strconv.ParseFloat(i, 64)
		if err!=nil {
			log.Printf("calcTotalPrice: %s", err)
			return totalPrice, err

		}
	}
	return totalPrice, err
}

//Calculates total distance of the flightplan
func calcTotalDistance(data []RouteInfo) int {
	var totalDistance int
	for _, flight := range data {
		totalDistance += flight.Distance

	}
	return totalDistance
}


//Finds the longesst duration
func findMaximumDuration(flights []Flights) int {
	var max int
	for i, flight := range flights {
		if i == 0 || int(flight.TravelTime.Hours()) > max {
			max = int(flight.TravelTime.Hours())
		}
	}
	return max+1
}
