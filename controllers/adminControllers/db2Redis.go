package adminControllers

import (
	"../../db/redis"
	"../../db/cassandra"
	"../../db/mysql"
	"net/http"
	"fmt"
	"strings"
	"strconv"
	"time"	
)

func UpdateHotelName2Redis(w http.ResponseWriter, r *http.Request) {
	mysql := mysql.Connect()
	rows, err := mysql.Query("SELECT hotel_id, hotel_name, city, country, rating_average FROM hotels_rank ORDER BY rating_average DESC, star_rating DESC, number_of_reviews DESC LIMIT 363617, 500000")
	// rows, err := mysql.Query("SELECT CONCAT(hotel_id, ':', hotel_name ,', ', city, ', ', country, ':', rating_average) AS v FROM hotels_rank ORDER BY rating_average DESC, star_rating DESC, number_of_reviews DESC LIMIT 25")
	if err != nil {
		fmt.Println(err)
	}
	// defer rows.Close()	
	rank:=0
	i:=0
	no:=0
	var v string
	var hotel_id, hotel_name, city, country, rating_average string
	fmt.Println("Hotels to Redis Starting...")
	for rows.Next() {
		rank=rank+100
		// err = rows.Scan(&v)
		err = rows.Scan(&hotel_id, &hotel_name, &city, &country, &rating_average)
		v = hotel_id +":"+ hotel_name +", "+ city +", "+ country +":"+ rating_average
		// fmt.Println(v)
		if err != nil {
			fmt.Print(err)
		}

		redis.ZAdd("hotels", rank, v)

		i++
		if i >= 500 {
			i = 0
			no++
			fmt.Println(no, ".\t500 rows added to Redis, pause 1 seconds....")
			time.Sleep(1*time.Second)
			fmt.Println("Continue...")
		}		
	}	
}




type City struct {
	City_id string
	Count int
	City string
	Country string
	Countryisocode string
}
func UpdateCities2Redis(w http.ResponseWriter, r *http.Request) {
	c := City{}
	fmt.Println("cities...")
	no:=0
	i:=0
	iter := cassandra.Session.Query("SELECT city_id, count, city, country, countryisocode FROM cities LIMIT 22210, 30000").Iter()
	for iter.Scan(&c.City_id, &c.Count, &c.City, &c.Country, &c.Countryisocode) {
		// v := strings.Join([]string{c.City_id, ":", c.City, ":", c.City, ":", c.Country, ":", c.Countryisocode}, "")
		v := strings.Join([]string{c.City_id, ":", c.City, ", ", c.Country, ":", strconv.Itoa(c.Count)}, "")
		// fmt.Println(v)
		redis.ZAdd("cities", c.Count, v)

		i++
		if i >= 500 {
			i = 0
			no++
			fmt.Println(no, ".\t500 rows added to Redis, pause 1 seconds....")
			time.Sleep(1*time.Second)
			fmt.Println("Continue...")
		}
	}	
}