package controllers

import (
	"../db/cassandra"
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

func Cities(w http.ResponseWriter, r *http.Request) {

	url := strings.ToLower(r.URL.String())
	url_replacer := strings.NewReplacer(".htm", "", "/", "")
	getCountrykey := url_replacer.Replace(url)
	fmt.Println(getCountrykey)

	cassandraSession := cassandra.Session

	type typeCountry struct {
		Country_id     string
		Count          string
		Country        string
		Countryisocode string
		Countrykey     string
	}

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
		Title     string
		Cities    [16]typeCity
		Countries [16]typeCountry
	}

	var datas = data{}

	//////////CITY	
	city_ids := []string{"17193", "5085", "9395", "233", "16056", "15470", "14690", "14552", "16594", "1569", "8584", "13170", "3987", "14932", "7401", "14524"}
	city_ids_in := strings.Join(city_ids, ",")
	cy := typeCity{}

	if getCountrykey != "" {
		iter := cassandraSession.Query("SELECT city_id, city, citykey, count, country, country_id, countryisocode, countryisocode_lower, countrykey FROM cities_by_country WHERE countrykey = '" + getCountrykey + "' ORDER BY count DESC LIMIT 16").Iter()
		i := 0
		for iter.Scan(&cy.City_id, &cy.City, &cy.Citykey, &cy.Count, &cy.Country, &cy.Country_id, &cy.Countryisocode, &cy.Countryisocode_lower, &cy.Countrykey) {
			datas.Cities[i] = typeCity{City_id: cy.City_id, City: cy.City, Citykey: cy.Citykey, Count: cy.Count, Country: cy.Country, Country_id: cy.Country_id, Countryisocode: cy.Countryisocode, Countryisocode_lower: cy.Countryisocode_lower, Countrykey: cy.Countrykey}
			i++
		}
	} else {
		iter := cassandraSession.Query("SELECT city_id, city, citykey, count, country, country_id, countryisocode, countryisocode_lower, countrykey FROM cities WHERE city_id IN ("+ city_ids_in +")").Iter()

		cities := map[string]typeCity{}
		for iter.Scan(&cy.City_id, &cy.City, &cy.Citykey, &cy.Count, &cy.Country, &cy.Country_id, &cy.Countryisocode, &cy.Countryisocode_lower, &cy.Countrykey) {
			cities[(cy.City_id)] = typeCity{City_id: cy.City_id, City: cy.City, Citykey: cy.Citykey, Count: cy.Count, Country: cy.Country, Country_id: cy.Country_id, Countryisocode: cy.Countryisocode, Countryisocode_lower: cy.Countryisocode_lower, Countrykey: cy.Countrykey}
		}
		for i := 0; i < len(cities); i++ {
			datas.Cities[i] = cities[(city_ids[i])]
		}

		if getCountrykey != "" {
			datas.Title = "Top Cities in " + datas.Cities[0].Country
		} else {
			datas.Title = "Top Cities Worldwide"
		}
	}


	////////COUNTRY
	country_ids := []string{"181", "191", "106", "35", "3", "192", "139", "38", "107", "140", "198", "212", "153", "70", "205", "167"}
	country_ids_in := strings.Join(country_ids, ",")

	iter2 := cassandraSession.Query("SELECT country_id, count, country, countryisocode, countrykey FROM countries WHERE country_id IN ("+ country_ids_in +")").Iter()

	countries := map[string]typeCountry{}
	co := typeCountry{}
	for iter2.Scan(&co.Country_id, &co.Count, &co.Country, &co.Countryisocode, &co.Countrykey) {
		countries[(co.Country_id)] = typeCountry{Country_id: co.Country_id, Count: co.Count, Country: co.Country, Countryisocode: co.Countryisocode, Countrykey: co.Countrykey}
	}
	for i := 0; i < len(countries); i++ {
		datas.Countries[i] = countries[(country_ids[i])]
	}	

	var t, err = template.ParseFiles("views/home.html")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	t.Execute(w, datas)
}
