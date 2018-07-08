package adminControllers

import (
	"../../db/cassandra"
	"fmt"
	"html/template"
	"net/http"
)

func Countries(w http.ResponseWriter, r *http.Request) {

	session := cassandra.Session

	type typeCountry struct {
		Country_id     string
		Count          string
		Country        string
		Countryisocode string
		Countrykey     string
	}

	type data struct {
		Wherestr  template.HTML
		Countries [15]typeCountry
	}

	var datas = data{}

	var search_label = "Results for"
	query_where := ""
	filter := false
	i := 0
	if country, ok := r.URL.Query()["country"]; ok {
		if country[0] != "" {
			i++
			if i == 1 {
				query_where = " WHERE country = '" + country[0] + "'"
			}
			search_label = search_label + " <b>country</b> = " + country[0]
			filter = true
		}
	}

	query := "SELECT * FROM country_oby_count" + query_where + " LIMIT 15 ALLOW FILTERING"
	fmt.Println(query)
	if filter == true {
		datas.Wherestr = template.HTML(search_label)
	}

	var country_id, count, country, countryisocode, countrykey string
	j := 0

	iter := session.Query(query).Iter()

	for iter.Scan(&country_id, &count, &country, &countryisocode, &countrykey) {
		datas.Countries[j] = typeCountry{Country_id: country_id, Count: count, Country: country, Countryisocode: countryisocode, Countrykey: countrykey}
		j++
	}

	var t, err = template.ParseFiles("views/admin/countries.html")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	t.Execute(w, datas)
}
