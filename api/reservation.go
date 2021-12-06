package api

import (
	"database/sql"
	"log"
	"net/url"
	"strconv"
	"time"

	"github.com/lib/pq"
)

type Reservation struct {
	ID          string
	PriceListID string
	FirstName   string
	LastName    string
	TotalPrice  float64
	TotalTime   string
	FlightId    []string
}

//Function to send reservation data to the database
func MakeReservation(form url.Values) error {
	r, err := ParseFormData(form)

	if err != nil {
		return err
	}

	db, err := sql.Open("postgres", connectionString)

	if err != nil {
		log.Printf("MakeReservation: %s", err)
		return err
	}

	defer db.Close()

	sqlStatement := `
	INSERT INTO reservation (id, pricelist_id, first_name, last_name, total_price, total_time, flights)
	SELECT CAST($1 AS VARCHAR), $2, $3, $4, $5, $6, $7
	WHERE NOT EXISTS (
	SELECT 1 FROM reservation WHERE id=$1
	);`
	_, err = db.Exec(sqlStatement, r.ID, r.PriceListID, r.FirstName, r.LastName, r.TotalPrice, r.TotalTime, pq.Array(r.FlightId))
	if err != nil {
		log.Printf("MakeReservation: %s", err)
		return err
	}
	return err
}


//Parses form data to Reservation struct
func ParseFormData(form url.Values) (Reservation, error) {

	var reservation Reservation
	stringTotalPrice := form.Get("totalprice")
	totalPrice, err := strconv.ParseFloat(stringTotalPrice, 64)

	if err != nil {
		log.Printf("ParseFormData: %s", err)
		return reservation, err
	}

	var flights []string
	for _, v := range form["flightid"] {
		flights = append(flights, v)
	}

	reservation = Reservation{
		form.Get("token"),
		form.Get("pricelistid"),
		form.Get("firstname"),
		form.Get("lastname"),
		totalPrice,
		form.Get("totaltime"),
		flights,
	}

	return reservation, err
}

//Checks, if the pricelist is still valid
func IsTimeValid(t string) (bool, error) {
	layout := "2006-01-02 15:04:05.999999 +0000 MST"
	validTime, err := time.Parse(layout, t)

	if err != nil {
		log.Printf("ParseFormData: %s", err)
		return false, err
	}

	if validTime.After(time.Now().UTC()) {
		return true, err
	}

	return false, err
}
