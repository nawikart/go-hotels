package cassandraModels

import (
	"../../db/cassandra"
	// "fmt"
	"strings"
)

type Hotel struct {
	Hotel_id             string
	Hotel_name           string
	City                 string
	Citykey              string
	Country              string
	Countryisocode       string
	Countryisocode_lower string
	Brand                string
	Addressline1         string
	Star_rating          string
	Photo                string
	// Overview             string
	No_reviews    string
	Hotel_namekey string
}

type Nearby struct {
	Hotels       map[int]Hotel
	Autocomplete string
}

func NearbyHotels(autocomplete string, hotel_ids []string) Nearby {

	hotel_ids_in := strings.Join(hotel_ids, ",")
	cassandraSession := cassandra.Session
	// defer CassandraSession.Close()
	query := "SELECT hotel_id, hotel_name, city, citykey, country, countryisocode, brand_name, addressline1, star_rating, photo1, number_of_reviews, hotel_namekey FROM hotels WHERE hotel_id IN (" + hotel_ids_in + ")"

	// fmt.Println(query)

	iter := cassandraSession.Query(query).Iter()
	hotel := Hotel{}
	hotels := map[string]Hotel{}
	hotels_sorted := Nearby{}
	hotels_sorted.Autocomplete = autocomplete
	hotels_sorted.Hotels = make(map[int]Hotel) //inisialisasi type map!

	for iter.Scan(&hotel.Hotel_id, &hotel.Hotel_name, &hotel.City, &hotel.Citykey, &hotel.Country, &hotel.Countryisocode, &hotel.Brand, &hotel.Addressline1, &hotel.Star_rating, &hotel.Photo, &hotel.No_reviews, &hotel.Hotel_namekey) {
		hotels[(hotel.Hotel_id)] = Hotel{Hotel_id: hotel.Hotel_id, Hotel_name: hotel.Hotel_name, City: hotel.City, Citykey: hotel.Citykey, Country: hotel.Country, Countryisocode: hotel.Countryisocode, Countryisocode_lower: strings.ToLower(hotel.Countryisocode), Brand: hotel.Brand, Addressline1: hotel.Addressline1, Star_rating: hotel.Star_rating, Photo: hotel.Photo, No_reviews: hotel.No_reviews, Hotel_namekey: hotel.Hotel_namekey}
	}

	for i := 0; i < len(hotels); i++ {
		hotels_sorted.Hotels[i] = hotels[(hotel_ids[i])]
	}

	// CassandraSession.Close()
	return hotels_sorted
}
