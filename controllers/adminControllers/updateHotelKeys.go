package adminControllers

import (
	"../../db/cassandra"
	"../../custom"
	"fmt"
	"time"
	"net/http"
)

func UpdateHotelKeys(w http.ResponseWriter, r *http.Request) {

	//PASTIKAN INI SUDAH
	// ALTER TABLE hotel ADD hotel_namekey text;
	// CREATE INDEX ON hotel (hotel_namekey);

	// ALTER TABLE hotel ADD citykey text;
	// CREATE INDEX ON hotel (citykey);

	// ALTER TABLE hotel ADD countrykey text;
	// CREATE INDEX ON hotel (countrykey);
	///////////////////

	session := cassandra.Session

	fmt.Println("STARTING....")

	var cql, hotel_id, hotel_name, city, country string
	iter := session.Query("SELECT hotel_id, hotel_name, city, country FROM hotels").Iter()
	i := 0
	j := 0
	var startTime = time.Now()
	for iter.Scan(&hotel_id, &hotel_name, &city, &country) {
		cql = "UPDATE hotels SET hotel_namekey = '" + custom.String2key(hotel_name) + "', citykey = '" + custom.String2key(city) + "' WHERE hotel_id = " + hotel_id
		// fmt.Println(cql)
		session.Query(cql).Iter()
		i++
		if(i >= 1000){
			j++
			i = 0
			fmt.Println(j, "). 1000 rows updated in", time.Since(startTime), "second")
			startTime = time.Now()
		}
	}

	fmt.Fprintln(w, "update 'hotel_namekey' 'citykey' 'countrykey' finish less than 10mnt!")
}
