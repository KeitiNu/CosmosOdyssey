package api

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

var (
	port             = "5432"
	user             = "postgres"
	dbname           = "odyssey"
	password         = "psw"
	connectionString = fmt.Sprintf("host=localhost port=%v user=%v dbname=%v sslmode=disable password=%v", port, user, dbname, password)
)

//Saves data from api to the database
func apiInfoToDB() error {

	db, err := sql.Open("postgres", connectionString)

	if err != nil {
		log.Printf("apiInfoToDB: %s", err)
		return err
	}

	defer db.Close()

	sqlStatement := `SELECT id FROM pricelist WHERE id=$1;`
	var id string
	row := db.QueryRow(sqlStatement, rawData.ID)

	switch err := row.Scan(&id); err {
	case sql.ErrNoRows:
		sqlStatement = `
		INSERT INTO pricelist (id, valid_until, created)
		SELECT CAST($1 AS VARCHAR), $2, $3
		WHERE NOT EXISTS (
		SELECT 1 FROM pricelist WHERE id=$1
		);`
		_, err = db.Exec(sqlStatement, rawData.ID, rawData.ValidUntil, time.Now())
		if err != nil {
			log.Printf("apiInfoToDB: %s", err)
			return err
		}

		for _, v := range rawData.Legs {
			sqlStatement := `
			INSERT INTO route_info (id, pricelist_id, from_location, to_location, distance)
			SELECT CAST($1 AS VARCHAR), $2, $3, $4, $5
			WHERE NOT EXISTS (
			SELECT 1 FROM route_info WHERE id=$1
			);`
			_, err = db.Exec(sqlStatement, v.RouteInfo.ID, rawData.ID, v.RouteInfo.From.Name, v.RouteInfo.To.Name, v.RouteInfo.Distance)
			if err != nil {
				log.Printf("apiInfoToDB: %s", err)
				return err
			}

			for _, flight := range v.Providers {

				sqlStatement := `
				INSERT INTO company (id, name, pricelist_id)
				SELECT CAST($1 AS VARCHAR), $2, $3
				WHERE NOT EXISTS (
				SELECT 1 FROM company WHERE id=$1
				);`
				_, err = db.Exec(sqlStatement, flight.Company.ID, flight.Company.Name, rawData.ID)
				if err != nil {
					log.Printf("apiInfoToDB: %s", err)
					return err
				}
			}

			for _, flight := range v.Providers {
				sqlStatement := `
				INSERT INTO flight (id, route_info_id, company_id, price, start_time, end_time)
				SELECT CAST($1 AS VARCHAR), $2, $3, $4, $5, $6
				WHERE NOT EXISTS (
				SELECT 1 FROM flight WHERE id=$1
				);`
				_, err = db.Exec(sqlStatement, flight.ID, v.RouteInfo.ID, flight.Company.ID, flight.Price, flight.FlightStart, flight.FlightEnd)
				if err != nil {
					log.Printf("apiInfoToDB: %s", err)
					return err
				}
			}

		}
	case nil:
	default:
		return err
	}

	sqlStatement = `DELETE FROM  pricelist WHERE
	ID not in (SELECT id FROM pricelist ORDER BY created DESC LIMIT 15)`

	_, err = db.Exec(sqlStatement)
	if err != nil {
		log.Printf("apiInfoToDB: %s", err)
		return err
	}

	return err
}
