package controllers

import (
	"../models/cassandraModels"
	"../models/mysqlModels"
	"fmt"
	"html/template"
	"net/http"
)

func NearbyHotels(w http.ResponseWriter, r *http.Request) {

	get := r.URL.Query()
	lat := get["lat"][0]
	lng := get["lng"][0]
	autocomplete := get["autocomplete"][0]

	var t, err = template.ParseFiles("views/nearbyhotels.html")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	t.Execute(w, cassandraModels.NearbyHotels(autocomplete, mysqlModels.NearbyHotelIDs(lat, lng)))
}
