package adminControllers

import (
	"../../db/cassandra"
	"fmt"
	"html/template"
	"net/http"
)

func Hotels(w http.ResponseWriter, r *http.Request) {

	session := cassandra.Session

	type hotel struct {
		Hotel_id     string
		Hotel_name   string
		City         string
		Country      string
		Brand        string
		Location     string
		Addressline1 string
		Star_rating  string
		Photo        string
		Overview     string
		No_reviews   string
	}

	type data struct {
		Wherestr template.HTML
		Hotels   [15]hotel
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
	if city, ok := r.URL.Query()["city"]; ok {
		if city[0] != "" {
			i++
			if i > 1 {
				search_label = search_label + ", "
				query_where = query_where + " AND city = '" + city[0] + "'"
			} else {
				query_where = " WHERE city = '" + city[0] + "'"
			}
			search_label = search_label + " <b>city</b> = " + city[0]
			filter = true
		}
	}
	if area, ok := r.URL.Query()["area"]; ok {
		if area[0] != "" {
			i++
			if i > 1 {
				search_label = search_label + ", "
				query_where = query_where + " AND area = '" + area[0] + "'"
			} else {
				query_where = " WHERE area = '" + area[0] + "'"
			}
			search_label = search_label + " <b>area</b> = " + area[0]
			filter = true
		}
	}
	if rank, ok := r.URL.Query()["rank"]; ok {
		if rank[0] != "" {
			i++
			if i > 1 {
				search_label = search_label + ", "
				query_where = query_where + " AND rank = '" + rank[0] + "'"
			} else {
				query_where = " WHERE rank = '" + rank[0] + "'"
			}
			search_label = search_label + " <b>rank</b> = " + rank[0]
			filter = true
		}
	}
	if star, ok := r.URL.Query()["star"]; ok {
		if star[0] != "" {
			i++
			if i > 1 {
				search_label = search_label + ", "
				query_where = query_where + " AND star = '" + star[0] + "'"
			} else {
				query_where = " WHERE star = '" + star[0] + "'"
			}
			search_label = search_label + " <b>star</b> = " + star[0]
			filter = true
		}
	}

	query := "SELECT hotel_id, hotel_name, city, country, brand_name, addressline1, star_rating, photo1, overview, number_of_reviews FROM hotels" + query_where + " LIMIT 15 ALLOW FILTERING"

	if filter == true {
		datas.Wherestr = template.HTML(search_label)
	}

	var hotel_id, hotel_name, city, country, brand_name, addressline1, star_rating, photo1, overview, number_of_reviews string
	j := 0

	iter := session.Query(query).Iter()

	for iter.Scan(&hotel_id, &hotel_name, &city, &country, &brand_name, &addressline1, &star_rating, &photo1, &overview, &number_of_reviews) {
		datas.Hotels[j] = hotel{Hotel_name: hotel_name, Hotel_id: hotel_id, City: city, Country: country, Brand: brand_name, Addressline1: addressline1, Star_rating: star_rating, Photo: photo1, Overview: overview, No_reviews: number_of_reviews}
		j++
	}

	var t, err = template.ParseFiles("views/admin/hotels.html")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	t.Execute(w, datas)
}
