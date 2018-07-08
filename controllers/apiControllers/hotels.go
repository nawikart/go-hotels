package apiControllers

import (
	"github.com/gocql/gocql"
	"../../db/mysql"
	"strings"
	"fmt"
	"strconv"
	"encoding/json"	
)

func hotels(dataFilter string, hotel_ids string, getCitykey string, getCountryIsoCodeLower string, page string) interface{}{

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
		CityId string
		Citykey string
		Country string
		Countrykey string
		Countryisocode_lower string
		P int
		Ps int
		Hotels []Hotel
	}

	var query string
	var total int

	if(dataFilter != ""){
		// fmt.Println(dataFilter)
		dataFilter_byt := []byte(dataFilter)
		var dataFilter_map map[string]string
		
		if err := json.Unmarshal(dataFilter_byt, &dataFilter_map); err != nil {
			panic(err)
		}
		
		var sql_dailyRate_min, sql_dailyRate_max, sql_minReviewScore, sql_minStarRating string
		sql_sortby := " ORDER BY rating_average DESC, star_rating DESC, number_of_reviews DESC"
		if(dataFilter_map["sortBy"] != ""){
			switch dataFilter_map["sortBy"] {
			case "PriceDesc":
				sql_sortby = " ORDER BY rates_from DESC, rating_average DESC, star_rating DESC, number_of_reviews DESC"
			case "PriceAsc":
				sql_sortby = " ORDER BY rates_from ASC, rating_average DESC, star_rating DESC, number_of_reviews DESC"
			case "StarRatingDesc":
				sql_sortby = " ORDER BY star_rating DESC, rating_average DESC, number_of_reviews DESC"
			case "StarRatingAsc":
				sql_sortby = " ORDER BY star_rating ASC, rating_average DESC, number_of_reviews DESC"
			case "ReviewScoreDesc":
				sql_sortby = " ORDER BY rating_average DESC, star_rating DESC, number_of_reviews DESC"
			case "ReviewScoreAsc":
				sql_sortby = " ORDER BY rating_average ASC, star_rating DESC, number_of_reviews DESC"
			case "ReviewsCountDesc":
				sql_sortby = " ORDER BY number_of_reviews DESC, rating_average DESC, star_rating DESC"
			case "ReviewsCountAsc":
				sql_sortby = " ORDER BY number_of_reviews ASC, rating_average DESC, star_rating DESC"
			}			
		}
		if(dataFilter_map["dailyRate_min"] != ""){
			dailyRate_min := strings.TrimLeft(dataFilter_map["dailyRate_min"], "a")
			sql_dailyRate_min = " AND (rates_from > "+ dailyRate_min +" OR rates_from = "+ dailyRate_min +")"
		}
		if(dataFilter_map["dailyRate_max"] != ""){
			dailyRate_max := strings.TrimLeft(dataFilter_map["dailyRate_max"], "a")
			sql_dailyRate_max = " AND (rates_from < "+ dailyRate_max +" OR rates_from = "+ dailyRate_max +")"
		}
		if(dataFilter_map["minReviewScore"] != ""){
			minReviewScore := strings.TrimLeft(dataFilter_map["minReviewScore"], "a")
			sql_minReviewScore = " AND (rating_average > "+ minReviewScore +" OR rating_average = "+ minReviewScore +")"
		}
		if(dataFilter_map["minStarRating"] != ""){
			minStarRating := strings.TrimLeft(dataFilter_map["minStarRating"], "a")
			sql_minStarRating = " AND (star_rating > "+ minStarRating +" OR star_rating = "+ minStarRating +")"
		}				
		
		cityId := strings.TrimLeft(dataFilter_map["cityId"], "a")
		// fmt.Println("cityId: ", cityId)
		sql_filter := "SELECT hotel_id FROM hotels_filter WHERE city_id = "+ cityId + sql_dailyRate_min + sql_dailyRate_max + sql_minReviewScore + sql_minStarRating + sql_sortby +" LIMIT 0, 24"
		fmt.Println(sql_filter)
		
		mysql := mysql.Connect()
		rows, err := mysql.Query(sql_filter)
		if err != nil {
			fmt.Println(err)
		}
		var hotel_ids_2 []string
		var hotel_id string
		for rows.Next() {
			err = rows.Scan(&hotel_id)
			if err != nil {
				fmt.Print(err)
			}
			// fmt.Println(hotel_id)
			if(hotel_id != ""){
				hotel_ids_2 = append(hotel_ids_2, hotel_id)
			}
		}
		if(len(hotel_ids_2) > 0){
			hotel_ids = strings.Join(hotel_ids_2,",")
		}
	}
	
	session.Query("SELECT count FROM cities WHERE citykey = '" + getCitykey + "' ALLOW FILTERING").Iter().Scan(&total)

	if hotel_ids != ""{
		query = "SELECT hotel_id, hotel_name, city_id, city, country, countryisocode, brand_name, addressline1, star_rating, photo1, overview, number_of_reviews, hotel_namekey, citykey, rates_from, rating_average FROM hotels WHERE hotel_id IN (" + hotel_ids + ") LIMIT 24 ALLOW FILTERING"
	}else if getCitykey != "" && getCountryIsoCodeLower != "" {		
		session.Query("SELECT hotel_ids FROM pagination_per_city WHERE citykey = '" + getCitykey + "' AND countryisocode_lower = '" + getCountryIsoCodeLower + "' AND p = "+ page +" ALLOW FILTERING").Iter().Scan(&hotel_ids)		
		query = "SELECT hotel_id, hotel_name, city_id, city, country, countryisocode, brand_name, addressline1, star_rating, photo1, overview, number_of_reviews, hotel_namekey, citykey, rates_from, rating_average FROM hotels WHERE hotel_id IN (" + hotel_ids + ") LIMIT 24 ALLOW FILTERING"
	}	

	var hotel_id, hotel_name, city_id, city, country, countryisocode, brand_name, addressline1, star_rating, photo1, overview, number_of_reviews, hotel_namekey, citykey, rates_from, rating_average string

	iter := session.Query(query).Iter()
	hotels := map[string]Hotel{}
	hotels_sorted := Data{}
	for iter.Scan(&hotel_id, &hotel_name, &city_id, &city, &country, &countryisocode, &brand_name, &addressline1, &star_rating, &photo1, &overview, &number_of_reviews, &hotel_namekey, &citykey, &rates_from, &rating_average) {
		hotels[hotel_id] = Hotel{Hotel_name: hotel_name, Hotel_namekey: hotel_namekey, Hotel_id: hotel_id, City: city, Citykey: citykey, Country: country, Countryisocode: countryisocode, Countryisocode_lower: strings.ToLower(countryisocode), Brand: brand_name, Addressline1: addressline1, Star_rating: star_rating, Photo: strings.Replace(photo1, "312x", "600x", 1), Overview: overview, No_reviews: number_of_reviews, Rating_average: rating_average, Rates_from: rates_from}
	}

	hotel_ids_arr := strings.Split(hotel_ids, ",")
	for i := 0; i < len(hotels); i++ {
		if hotels[(hotel_ids_arr[i])].Hotel_id != "" {
			hotels_sorted.Hotels = append(hotels_sorted.Hotels, hotels[(hotel_ids_arr[i])])
		}
	}

	hotels_sorted.City = city
	hotels_sorted.CityId = city_id
	hotels_sorted.Citykey = getCitykey
	hotels_sorted.Country = country
	hotels_sorted.Countryisocode_lower = getCountryIsoCodeLower
	hotels_sorted.P, _ = strconv.Atoi(page)
	hotels_sorted.Ps = total/24
	
	if((total%24)>0){hotels_sorted.Ps++}

	return hotels_sorted
}