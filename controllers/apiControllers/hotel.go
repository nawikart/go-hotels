package apiControllers

import (
	"../../db/cassandra"
	// "fmt"
	"html/template"
	// "net/http"
	// "regexp"
	"strings"
	"../../custom"
)

func hotel(where string) interface{}{

	cassandraSession := cassandra.Session

	type data struct {
		Wherestr             template.HTML
		Hotel_id             string
		Hotel_name           string
		Brand_name           string
		City                 string
		Citykey              string
		Country              string
		Countrykey              string
		Countryisocode       string
		Countryisocode_lower string
		Brand                string
		Location             string
		Addressline1         string
		Star_rating          string
		Photo1               string
		Photo2               string
		Photo3               string
		Photo4               string
		Overview             string
		No_reviews           string
		Hotel_namekey        string
		Number_of_reviews    string
		Rates_from string
		Rating_average string
	}

	var datas = data{}
	datas.Wherestr = "Ccccc"

	query := "SELECT hotel_id, hotel_name, city, country, countryisocode, brand_name, addressline1, star_rating, photo1, photo2, photo3, photo4, overview, number_of_reviews, rates_from, rating_average FROM hotels "+ where +" LIMIT 1 ALLOW FILTERING"
	// fmt.Println(query)
	var hotel_id, hotel_name, city, country, countryisocode, brand_name, addressline1, star_rating, photo1, photo2, photo3, photo4, overview, number_of_reviews, rates_from, rating_average string

	cassandraSession.Query(query).Iter().Scan(&hotel_id, &hotel_name, &city, &country, &countryisocode, &brand_name, &addressline1, &star_rating, &photo1, &photo2, &photo3, &photo4, &overview, &number_of_reviews, &rates_from, &rating_average)
	datas.Hotel_id = hotel_id
	datas.Hotel_name = hotel_name
	datas.City = city
	datas.Country = country
	datas.Countryisocode = countryisocode
	datas.Brand_name = brand_name
	datas.Addressline1 = addressline1
	datas.Star_rating = star_rating
	datas.Photo1 = strings.TrimRight(photo1, "?s=312x")
	datas.Photo2 = strings.TrimRight(photo2, "?s=312x")
	datas.Photo3 = strings.TrimRight(photo3, "?s=312x")
	datas.Photo4 = strings.TrimRight(photo4, "?s=312x")
	datas.Overview = overview
	datas.Number_of_reviews = number_of_reviews
	datas.Rates_from = rates_from
	datas.Rating_average = rating_average

	datas.Citykey = custom.String2key(city)
	datas.Countrykey = custom.String2key(country)
	datas.Countryisocode_lower = strings.ToLower(countryisocode)
	datas.Hotel_namekey = custom.String2key(hotel_name)

	return datas
}
