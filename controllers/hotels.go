package controllers

import (
	"../../github.com/gocql/gocql"
	"fmt"
	"html/template"
	"net/http"
	"regexp"
	"strings"
	"strconv"
)

func Hotels(w http.ResponseWriter, r *http.Request) {

	url := strings.ToLower(r.URL.String())
	rg, _ := regexp.Compile("/([-a-z]+)-([a-z][a-z])-hotels((-p([0-9]+))*).html")
	urlKeys := rg.FindStringSubmatch(url)
	// fmt.Println(urlKeys)
	getCitykey := urlKeys[1]
	getCountryIsoCodeLower := urlKeys[2]

	var page string
	if urlKeys[5] != "" {
		page = urlKeys[5]
	}else{
		page = "1"
	}
	

	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "mydb"
	session, _ := cluster.CreateSession()

	type Hotel struct {
		Hotel_id             string
		Hotel_name           string
		City                 string
		Citykey              string
		Country              string
		Countryisocode       string
		Countryisocode_lower string
		Brand                string
		Location             string
		Addressline1         string
		Star_rating          string
		Photo                string
		Overview             string
		No_reviews           string
		Hotel_namekey        string
		Rates_from 		 string
		Rating_average 		 string
	}

	type Data struct {
		City   string
		Citykey string
		Countryisocode_lower string
		P int
		Ps int
		Hotels map[int]Hotel
	}

	var query, hotel_ids string
	var total int

	if getCitykey != "" && getCountryIsoCodeLower != "" {
		session.Query("SELECT count FROM cities WHERE citykey = '" + getCitykey + "' ALLOW FILTERING").Iter().Scan(&total)
		session.Query("SELECT hotel_ids FROM pagination_per_city WHERE citykey = '" + getCitykey + "' AND countryisocode_lower = '" + getCountryIsoCodeLower + "' AND p = "+ page +" ALLOW FILTERING").Iter().Scan(&hotel_ids)
		query = "SELECT hotel_id, hotel_name, city, country, countryisocode, brand_name, addressline1, star_rating, photo1, overview, number_of_reviews, hotel_namekey, citykey, rates_from, rating_average FROM hotels WHERE hotel_id IN (" + hotel_ids + ") LIMIT 24 ALLOW FILTERING"
		// fmt.Println(query)
	}

	var hotel_id, hotel_name, city, country, countryisocode, brand_name, addressline1, star_rating, photo1, overview, number_of_reviews, hotel_namekey, citykey, rates_from, rating_average string

	iter := session.Query(query).Iter()
	hotels := map[string]Hotel{}
	hotels_sorted := Data{}
	hotels_sorted.Hotels = make(map[int]Hotel) //inisialisasi type map!	
	for iter.Scan(&hotel_id, &hotel_name, &city, &country, &countryisocode, &brand_name, &addressline1, &star_rating, &photo1, &overview, &number_of_reviews, &hotel_namekey, &citykey, &rates_from, &rating_average) {
		hotels[hotel_id] = Hotel{Hotel_name: hotel_name, Hotel_namekey: hotel_namekey, Hotel_id: hotel_id, City: city, Citykey: citykey, Country: country, Countryisocode: countryisocode, Countryisocode_lower: strings.ToLower(countryisocode), Brand: brand_name, Addressline1: addressline1, Star_rating: star_rating, Photo: photo1, Overview: overview, No_reviews: number_of_reviews, Rating_average: rating_average, Rates_from: rates_from}
		
	}
	hotel_ids_arr := strings.Split(hotel_ids, ",")
	for i := 0; i < len(hotels); i++ {
		hotels_sorted.Hotels[i] = hotels[(hotel_ids_arr[i])]
	}
	hotels_sorted.City = city
	hotels_sorted.Citykey = getCitykey
	hotels_sorted.Countryisocode_lower = getCountryIsoCodeLower
	hotels_sorted.P, _ = strconv.Atoi(page)
	hotels_sorted.Ps = total/24
	
	if((total%24)>0){hotels_sorted.Ps++}

	var t, err = template.ParseFiles("views/hotels.html")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	t.Execute(w, hotels_sorted)
}
