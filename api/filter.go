package api

import (
	"net/url"
	"sort"
	"strconv"
)

//Filters data by users choices
func Filter(form url.Values) (Data, error) {
	data, err := GetData()
	if err != nil {
		return data, err
	}

	possibleroutes := checkRoutes(form)
	data.FlightPlans, err = makeFlightPlan(possibleroutes)
	if err != nil {
		return data, err
	}

	if form.Has("from") {
		data.FormInfo.DepartureLocation = form.Get("from")

	}
	if form.Has("to") {
		data.FormInfo.Destination = form.Get("to")
	}

	data.FormInfo.MaxPrice = findMostExpensive(data.FlightPlans)
	data.FormInfo.MaxDistance = findLongest(data.FlightPlans)
	data.FormInfo.MaxTime = findMaximumDuration(data.FlightPlans)

	if form.Has("providers") {
		data.UniqueProviders = checkProviders(form, data.UniqueProviders)
		data.FlightPlans = filterByProviders(data.FlightPlans, form["providers"])
	}

	if form.Has("price") {
		price, _ := strconv.ParseFloat(form.Get("price"), 64)
		data.FormInfo.ChosenPrice = int(price)
		data.FlightPlans = filterByPrice(data.FlightPlans, price)
	}

	if form.Has("distance") {
		data.FormInfo.ChosenDistance, _ = strconv.Atoi(form.Get("distance"))
		data.FlightPlans = filterByDistance(data.FlightPlans, data.FormInfo.ChosenDistance)
	}

	if form.Has("time") {
		data.FormInfo.ChosenTime, _ = strconv.Atoi(form.Get("time"))
		data.FlightPlans = filterByTime(data.FlightPlans, data.FormInfo.ChosenTime)
	}

	if form.Has("sort") {
		data.FormInfo.ChosenSort = form.Get("sort")
		data.FlightPlans = Sort(data.FlightPlans, data.FormInfo.ChosenSort)
	}

	return data, err
}

//Filters []Flights by providers
func filterByProviders(flights []Flights, validCompanies []string) []Flights {
	var filteredFlights []Flights

	for _, flight := range flights {
		valid := true
		for _, v := range flight.Flights { //Start comparing a company name vs user input
			match := true
			for _, s := range validCompanies {
				if v.Company.Name == s { //if it finds a match, break from the loop and compare other names, if all names get a match, then append the valid provider to the filtered list
					match = true
					break
				}
				match = false
			}
			if !match { //looped through valid companies, and found no match
				valid = false
				break //Breaks out of the loop, beacause it is pointless to compare other companies
			}
		}
		if valid { //if I get here, i want to know, if the flight was valid or not, if yes, then append, if no, nothing happens
			filteredFlights = append(filteredFlights, flight)
		}
	}
	return filteredFlights
}

//Filters []Flights by price
func filterByPrice(flights []Flights, maxPrice float64) []Flights {
	var filteredFlights []Flights

	for _, v := range flights {
		if v.TotalPrice <= maxPrice {
			filteredFlights = append(filteredFlights, v)
		}
	}

	return filteredFlights
}

//Filters []Flights by distance
func filterByDistance(flights []Flights, maxDistance int) []Flights {
	var filteredFlights []Flights

	for _, v := range flights {
		if v.TotalDistance <= maxDistance {
			filteredFlights = append(filteredFlights, v)
		}
	}

	return filteredFlights
}

//Filters []Flights by time
func filterByTime(flights []Flights, maxTime int) []Flights {
	var filteredFlights []Flights

	for _, v := range flights {
		if v.TravelTime.Hours() <= float64(maxTime)  {
			filteredFlights = append(filteredFlights, v)
		}
	}

	return filteredFlights
}

//Sorts []Flights by user input
func Sort(trips []Flights, p string) []Flights {
	sortedList := trips

	switch p {
	case "price":
		sort.Slice(sortedList, func(i, j int) bool {
			return int(sortedList[i].TotalPrice) < int(sortedList[j].TotalPrice)
		})

	case "distance":
		sort.Slice(sortedList, func(i, j int) bool {
			return sortedList[i].TotalDistance < sortedList[j].TotalDistance
		})

	case "duration":
		sort.Slice(sortedList, func(i, j int) bool {
			return int(sortedList[i].TravelTime) < int(sortedList[j].TravelTime)
		})
	default:
		sortedList = trips
	}

	return sortedList
}

//Marks the providers as checked or unchecked (as in form values) in a map
func checkProviders(form url.Values, providers map[string]bool) map[string]bool {

	//Uncheck all values
	for provider := range providers {
		providers[provider] = false 
	}

	//Check chosen values
	for _, v := range form["providers"] {
		providers[v] = true
	}

	return providers
}
