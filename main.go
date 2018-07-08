package main

import (
	"./controllers"
	"./controllers/adminControllers"
	"./controllers/apiControllers"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {

	r := mux.NewRouter()

	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/"))))

	///APi
	r.HandleFunc("/autocomplete/{key}", apiControllers.AutocompleteMysql)
	r.HandleFunc("/front", apiControllers.Front)

	r.HandleFunc("/autocomplete-redis/{key}", apiControllers.AutocompleteRedis)
	r.HandleFunc("/autocomplete-mongodb/{key}", apiControllers.AutocompleteMongodb)

	r.HandleFunc("/autocomplete-redirect", apiControllers.AutocompleteRedirect)

	/////FRONT
	r.HandleFunc("/nearby-hotels.html", controllers.NearbyHotels)
	r.HandleFunc("/{citykey}-{countryisocodelower}-hotels.html", controllers.Hotels)
	r.HandleFunc("/{citykey}-{countryisocodelower}-hotels{page:-p[0-9]+}.html", controllers.Hotels)
	r.HandleFunc("/{hotel_namekey}/hotel/{citykey}-{countryisocode_lower}.html", controllers.Hotel)
	r.HandleFunc("/{countrykey}.htm", controllers.Cities)
	r.HandleFunc("/", controllers.Cities)

	/////ADMIN
	r.HandleFunc("/admin", adminControllers.Home)
	r.HandleFunc("/admin/hotels/update-hotel-paging", adminControllers.UpdateHotelPaging)
	r.HandleFunc("/admin/hotels/update-hotel-keys", adminControllers.UpdateHotelKeys)
	// r.HandleFunc("/admin/hotels/update-hotel-hotel_namekey", adminControllers.UpdateHotelNameKeys)
	// r.HandleFunc("/admin/hotels/update-hotel-citykey", adminControllers.UpdateHotelCityKey)
	// r.HandleFunc("/admin/hotels/update-hotel-countrykey", adminControllers.UpdateHotelCountryKey)
	r.HandleFunc("/admin/hotels", adminControllers.Hotels)
	r.HandleFunc("/admin/countries", adminControllers.Countries)
	r.HandleFunc("/admin/cities", adminControllers.Cities)
	r.HandleFunc("/admin/cities/update-keys", adminControllers.UpdateCitiesKeys)

	r.HandleFunc("/admin/update-hotelname-to-redis", adminControllers.UpdateHotelName2Redis)
	r.HandleFunc("/admin/update-cities-to-redis", adminControllers.UpdateCities2Redis)

	http.Handle("/", r)

	fmt.Println("starting web server at http://localhost:8080/")
	http.ListenAndServe(":8080", nil)
}
