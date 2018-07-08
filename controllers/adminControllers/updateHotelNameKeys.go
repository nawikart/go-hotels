package adminControllers

import (
	"../../db/cassandra"
	"../../custom"
	"fmt"
	"net/http"
)

func UpdateHotelNameKeys(w http.ResponseWriter, r *http.Request) {

	//PASTIKAN INI SUDAH
	// ALTER TABLE hotel ADD hotel_namekey text;
	// CREATE INDEX ON hotel (hotel_namekey);
	///////////////////
	
	session := cassandra.Session

	var cql, hotel_id, hotel_name string
	iter := session.Query("SELECT hotel_id, hotel_name FROM hotel").Iter()
	for iter.Scan(&hotel_id, &hotel_name) {
		cql = "UPDATE hotel SET hotel_namekey = '" + custom.String2key(hotel_name) + "' WHERE hotel_id = " + hotel_id
		// fmt.Println(cql)
		session.Query(cql).Iter()
	}

	fmt.Fprintln(w, "update 'hotel_namekey' finish!")
}
