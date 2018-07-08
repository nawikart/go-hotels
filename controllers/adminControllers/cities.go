package adminControllers

import (
	"../../db/cassandra"
	"fmt"
	"html/template"
	"net/http"
)

func Cities(w http.ResponseWriter, r *http.Request) {

	session := cassandra.Session

	type typeCity struct {
		City_id              string
		City                 string
		Citykey              string
		Count                string
		Country              string
		Country_id           string
		Countryisocode       string
		Countryisocode_lower string
		Countrykey           string
	}

	type data struct {
		Wherestr template.HTML
		Cities   [15]typeCity
	}

	var datas = data{}

	var search_label = "Results for"
	query_where := ""
	filter := false
	i := 0
	if city, ok := r.URL.Query()["city"]; ok {
		if city[0] != "" {
			i++
			if i == 1 {
				query_where = " WHERE city = '" + city[0] + "'"
			}
			search_label = search_label + " <b>city</b> = " + city[0]
			filter = true
		}
	}

	query := "SELECT * FROM city_oby_count" + query_where + " LIMIT 15 ALLOW FILTERING"

	if filter == true {
		datas.Wherestr = template.HTML(search_label)
	}

	var city_id, city, citykey, count, country, country_id, countryisocode, countryisocode_lower, countrykey string
	j := 0

	iter := session.Query(query).Iter()

	for iter.Scan(&city_id, &city, &citykey, &count, &country, &country_id, &countryisocode, &countryisocode_lower, &countrykey) {
		datas.Cities[j] = typeCity{City_id: city_id, City: city, Citykey: citykey, Count: count, Country: country, Country_id: country_id, Countryisocode: countryisocode, Countryisocode_lower: countryisocode_lower, Countrykey: countrykey}
		j++
	}
	var t, err = template.ParseFiles("views/admin/cities.html")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	t.Execute(w, datas)
}
