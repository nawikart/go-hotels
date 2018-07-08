package adminControllers

import (
	"../../db/cassandra"
	"../../db/mysql"
	"net/http"
	"fmt"
	"strings"
	"strconv"
	// "time"
)

func UpdateHotelPaging(w http.ResponseWriter, r *http.Request) {

	mysql := mysql.Connect()
	var hotel_id, city_id, citykey, countryisocode string
	var hotel_ids = []string{}

	fmt.Println("cities...")
	limit:=24
	i:=0
	c:=0	
	p:=1
	no:=1
	iter := cassandra.Session.Query("SELECT city_id, citykey, countryisocode FROM cities WHERE city_id = 21740").Iter() // WHERE city_id = 17193
	for iter.Scan(&city_id, &citykey, &countryisocode) {
		p=1			
		rows, err := mysql.Query("SELECT hotel_id FROM hotels_rank WHERE city_id = "+ city_id +" ORDER BY rating_average DESC, star_rating DESC, number_of_reviews DESC")
		if err != nil {
			fmt.Println(err)
		}
		// defer rows.Close()
		i=0
		for rows.Next() {
			err = rows.Scan(&hotel_id)
			if err != nil {
				fmt.Print(err)
			}
			hotel_ids = append(hotel_ids, hotel_id)
			i++
			if(i>=limit){
				i=0
				cassandra.Session.Query("INSERT INTO pagination_per_city (city_id, citykey, countryisocode_lower, p, hotel_ids) VALUES ("+ city_id +", '"+ citykey +"', '"+ strings.ToLower(countryisocode) +"', "+ strconv.Itoa(p) +", '"+ strings.Join(hotel_ids, ",") +"')").Iter()
				hotel_ids=[]string{}
				p++
			}						
		}
		if(len(hotel_ids)>0){
			cassandra.Session.Query("INSERT INTO pagination_per_city (city_id, citykey, countryisocode_lower, p, hotel_ids) VALUES ("+ city_id +", '"+ citykey +"', '"+ strings.ToLower(countryisocode) +"', "+ strconv.Itoa(p) +", '"+ strings.Join(hotel_ids, ",") +"')").Iter()
			hotel_ids=[]string{}
		}

		no++
		c++
		if c >= 100 {
			c = 0
			fmt.Println(no, ".\t100 rows cities processed....")
			// fmt.Println(no, ".\t100 rows cities processing, pause 5 seconds....")
			// time.Sleep(5*time.Second)
			// mysql.Close()
			// fmt.Println("Continue...")
		}			
	}	
}