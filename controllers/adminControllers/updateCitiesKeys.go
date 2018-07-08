package adminControllers

import (
	"../../db/cassandra"
	"fmt"
	"../../custom"
	"net/http"
	"strings"
)

func UpdateCitiesKeys(w http.ResponseWriter, r *http.Request) {

	session := cassandra.Session

	var city_id, city, country, countryisocode string

	iter := session.Query("SELECT city_id, city, country, countryisocode FROM cities").Iter()

	for iter.Scan(&city_id, &city, &country, &countryisocode) {
		session.Query("UPDATE cities SET citykey = '" + custom.String2key(city) + "', countrykey = '" + custom.String2key(country) + "', countryisocode_lower = '" + strings.ToLower(countryisocode) + "'  WHERE city_id = " + city_id).Iter()
	}

	fmt.Fprintln(w, "update TABLE cities --> 'citykey, countrykey, countryisocode_lower' finish!")
}
