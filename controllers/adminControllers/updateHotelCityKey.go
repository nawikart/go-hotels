package adminControllers

import (
	"../../db/cassandra"
	"fmt"
	"net/http"
)

func UpdateHotelCityKey(w http.ResponseWriter, r *http.Request) {

	//BELUM SELESAI
	///////////////////

	//PASTIKAN INI SUDAH
	// ALTER TABLE hotel ADD citykey text;
	// CREATE INDEX ON hotel (citykey);
	///////////////////

	session := cassandra.Session

	var city_id, citykey string

	iter := session.Query("SELECT city_id, citykey FROM city_oby_count").Iter()

	for iter.Scan(&city_id, &citykey) {
		session.Query("UPDATE hotel SET citykey = '" + citykey + "' WHERE city_id = " + city_id).Iter()
	}

	fmt.Fprintln(w, "update 'citykey' finish!")
}
