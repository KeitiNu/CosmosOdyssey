package main

import (
	"cosmosodyssey/api"
	"fmt"
	"net/http"
	"text/template"
)

//Handles / htttp request
func indexHandler(w http.ResponseWriter, r *http.Request) {
	
	if r.Method == "GET" {

		if r.URL.Path != "/" {

			http.Error(w, "404: page not found", http.StatusNotFound)
			return
		}

		http.ServeFile(w, r, "../templates/index.html")

	} else {
		if r.URL.Path != "/" {

			http.Error(w, "404: page not found", http.StatusNotFound)
			return
		}

		http.Error(w, "405: method not allowed", http.StatusMethodNotAllowed)
		return

	}
}

//Handles /filter http request
func filterHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/filter" {

		http.Error(w, "404: page not found", http.StatusNotFound)
		return
	}

	if r.Method == "GET" {

		http.Redirect(w, r, "/", http.StatusSeeOther)

	} else if r.Method == "POST" {

		err := r.ParseForm()

		if err != nil {
			http.Error(w, "400: Bad Request", http.StatusBadRequest)
			return
		}

		t, err := template.ParseFiles("../templates/filter.html")

		if err != nil {
			http.Error(w, "500: internal server error", http.StatusInternalServerError)
			return
		}

		filteredData, err := api.Filter(r.Form)

		if err != nil {
			http.Error(w, "500: internal server error", http.StatusInternalServerError)
			return
		}

		t.Execute(w, filteredData)

	} else {
		http.Error(w, "405: method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

//Handles /reservation request
func reservationHandler(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/reservation" {
		http.Error(w, "404: page not found", http.StatusNotFound)
		return
	}

	if r.Method == "GET" {

		http.Error(w, "404: page not found", http.StatusNotFound)
		return

	} else if r.Method == "POST" {

		err := r.ParseForm()

		if err != nil {
			http.Error(w, "400: Bad Request", http.StatusBadRequest)
			return
		}

		t, err := template.ParseFiles("../templates/reservation.html")

		if err != nil {
			http.Error(w, "500: internal server error", http.StatusInternalServerError)
			return
		}

		valid, err := api.IsTimeValid(r.FormValue("valid"))
		if err != nil {
			http.Error(w, "500: internal server error", http.StatusInternalServerError)
			return
		}

		if valid { //Checks if the pricelist is still valid
			api.MakeReservation(r.Form)
			t.Execute(w, true)
		} else { //Cant make a reservation, because time is up
			t.Execute(w, false)
		}

	} else {
		http.Error(w, "405: method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func main() {

	http.Handle("/templates/", http.StripPrefix("/templates/", http.FileServer(http.Dir("../templates/"))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("../templates/static/css/"))))

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/filter", filterHandler)
	http.HandleFunc("/reservation", reservationHandler)

	fmt.Println("Server listening on port 3000")
	http.ListenAndServe(":3000", nil)

}
