package apiControllers

import (
	"strings"
	"../../custom"
	"../../db/redis"
	"../../db/cassandra"
	"../../db/mysql"
	"net/http"
	"encoding/json"
	"fmt"
)
type H struct {
	Hotel_id            string
	Hotel_name          string
}
type C struct {
	City_id             string
	City		        string
}
type S struct {
	Type	string
	Id      string
	Value   string
	Count   string
	Rating_average   string
}	
type A struct {
	Hotel_namekey string
	Citykey string
	CityId string
	Countryisocode string
}
func AutocompleteMysql(w http.ResponseWriter, r *http.Request) {

	url := strings.ToLower(r.URL.String())
	url_replacer := strings.NewReplacer("/autocomplete/", "")
	keyword := url_replacer.Replace(url)

	var se = S{}
	var s = []S{}
	mysql := mysql.Connect()
		
	rows, err := mysql.Query("SELECT city_id, CONCAT(city, ', ', country) AS v, count FROM cities_rank WHERE LOWER(city) LIKE '%"+ strings.ToLower(keyword) +"%' OR LOWER(country) LIKE '%"+ strings.ToLower(keyword) +"%' ORDER BY count DESC LIMIT 0, 5")
	if err != nil {
		fmt.Println(err)
	}
	for rows.Next() {
		err = rows.Scan(&se.Id, &se.Value, &se.Count)
		if err != nil {
			fmt.Print(err)
		}
		s = append(s, S{Type: "c", Id: se.Id, Value: se.Value, Count: se.Count})
	}


	rows, err = mysql.Query("SELECT hotel_id, CONCAT(hotel_name ,', ', city, ', ', country) AS v, rating_average FROM hotels_rank WHERE LOWER(hotel_name) LIKE '%"+ strings.ToLower(keyword) +"%' ORDER BY rating_average DESC, star_rating DESC, number_of_reviews DESC LIMIT 0, 5")
	if err != nil {
		fmt.Println(err)
	}
	for rows.Next() {
		err = rows.Scan(&se.Id, &se.Value, &se.Rating_average)
		if err != nil {
			fmt.Print(err)
		}
		s = append(s, S{Type: "h", Id: se.Id, Value: se.Value, Rating_average: se.Rating_average})
	}


	mysql.Close()

	js, err := json.Marshal(s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(js)
}

func autocompleteCityUrl(ty string, id string) A{
	a := A{}
	if(ty == "h"){
		cassandra.Session.Query("SELECT Hotel_namekey, citykey, city_id, countryisocode FROM hotels WHERE hotel_id = "+ id).Iter().Scan(&a.Hotel_namekey, &a.Citykey, &a.CityId, &a.Countryisocode)
	}else{
		cassandra.Session.Query("SELECT citykey, city_id, countryisocode FROM cities WHERE city_id = "+ id).Iter().Scan(&a.Citykey, &a.CityId, &a.Countryisocode)
	}

	return a
}

func AutocompleteRedis(w http.ResponseWriter, r *http.Request) {

	url := strings.ToLower(r.URL.String())
	url_replacer := strings.NewReplacer("/autocomplete-redis/", "")
	keyword := url_replacer.Replace(url)

	// fmt.Println(keyword)
	pattern := strings.Join([]string{"*", custom.String2InCaseSensitivePattern(keyword), "*"}, "")
	resH := redis.ZScan2("hotels", 0, pattern, 1000000)
	resC := redis.ZScan2("cities", 0, pattern, 30000)
	fmt.Println(pattern)
	
	var s = []S{}

	// var c = []C{}
	var vcs = []string{}
	i:=0
	for _, vc := range resC {
		vcs = strings.Split(vc, ":")
		if len(vcs) > 1 {
			s = append(s, S{Type: "c", Id: vcs[0], Value: vcs[1], Count: vcs[2]})
			i++
			if i > 5 {
				break
			}
		}
	}
	// var h = []H{}
	var vhs = []string{}
	i=0
    for _, vh := range resH {
		vhs = strings.Split(vh, ":")
		if len(vhs) > 1 {
			//h = append(h, H{Hotel_id: vs[0], Hotel_name: vs[1], Hotel_namekey: vs[2], City: vs[3], Citykey: vs[4], Country: vs[5], Countryisocode: vs[6]})
			s = append(s, S{Type: "h", Id: vhs[0], Value: vhs[1], Rating_average: vhs[2]})
			i++
			if i > 5 {
				break
			}			
		}
	}


	js, err := json.Marshal(s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
			return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(js)
}

func AutocompleteRedirect(w http.ResponseWriter, r *http.Request){
	a := A{}
	ty := r.URL.Query()["type"][0]
	id := r.URL.Query()["id"][0]
	var redirect string
	if(ty == "h"){
		cassandra.Session.Query("SELECT Hotel_namekey, citykey, countryisocode FROM hotels WHERE hotel_id = "+ id).Iter().Scan(&a.Hotel_namekey, &a.Citykey, &a.Countryisocode)
		redirect = a.Hotel_namekey +"/hotel/"+ a.Citykey +"-"+ strings.ToLower(a.Countryisocode) +".html"
	}else{
		cassandra.Session.Query("SELECT citykey, countryisocode FROM cities WHERE city_id = "+ id).Iter().Scan(&a.Citykey, &a.Countryisocode)
		redirect = a.Citykey +"-"+ strings.ToLower(a.Countryisocode) +"-hotels.html"
	}
	// fmt.Println(redirect)
	http.Redirect(w, r, redirect, 301)
}