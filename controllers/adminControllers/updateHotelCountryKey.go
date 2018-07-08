package adminControllers

import (
	"../../db/cassandra"
	"fmt"
	"net/http"
)

func UpdateHotelCountryKey(w http.ResponseWriter, r *http.Request) {

	//BELUM SELESAI
	///////////////////

	//PASTIKAN INI SUDAH
	// ALTER TABLE hotel ADD countrykey text;
	// CREATE INDEX ON hotel (countrykey);
	///////////////////

	session := cassandra.Session

	var country_id, countrykey string

	iter := session.Query("SELECT country_id, countrykey FROM country_oby_count").Iter()

	for iter.Scan(&country_id, &countrykey) {
		session.Query("UPDATE hotel SET countrykey = '" + countrykey + "' WHERE country_id = " + country_id).Iter()
	}

	fmt.Fprintln(w, "update 'countrykey' finish!")
}
