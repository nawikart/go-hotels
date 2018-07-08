package apiControllers

import (
	"encoding/json"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"strings"
	"../../custom"
	"strconv"
)

func connect() (*mgo.Session, error) {
	var session, err = mgo.Dial("localhost")
	if err != nil {
		return nil, err
	}
	return session, nil
}

type Sh struct {
	Hotel_id   int    `bson:"hotel_id"`	
	Hotel_name string `bson:"hotel_name"`
	City   string    `bson:"city"`
	Country   string    `bson:"country"`
	Rating_average   float64    `bson:"rating_average"`
}
type Sc struct {
    City_id int 
    City string
    Citykey string
    Count int
    Country_id int
    Country string
    Countrykey string
    Countryisocode string
    Countryisocode_lower string
}

func AutocompleteMongodb(w http.ResponseWriter, r *http.Request) {

	url := strings.ToLower(r.URL.String())
	url_replacer := strings.NewReplacer("/autocomplete-mongodb/", "")
	keyword := url_replacer.Replace(url)

	// fmt.Println(keyword)

	w.Header().Set("Content-Type", "application/json")

	var session, err = connect()
	if err != nil {
		fmt.Println("Error!", err.Error())
		return
	}
	defer session.Close()
	var collectionHotels = session.DB("test").C("hotels_rank")

	var hotels []Sh
	var selector = bson.M{"hotel_name": bson.M{"$regex": custom.String2InCaseSensitivePattern(keyword)}}
	err = collectionHotels.Find(selector).Sort("-rating_average").Limit(5).All(&hotels)
	if err != nil {
		fmt.Println("Error!", err.Error())
		return
	}

	var collectionCities = session.DB("test").C("cities_finish")
	var cities []Sc
	selector = bson.M{"city": bson.M{"$regex": custom.String2InCaseSensitivePattern(keyword)}}

	// selector = bson.M{ or:[{"city": bson.M{"$regex": custom.String2InCaseSensitivePattern(keyword)}}, {"country": bson.M{"$regex": custom.String2InCaseSensitivePattern(keyword)}}]}



	err = collectionCities.Find(selector).Sort("-count").Limit(5).All(&cities)
	if err != nil {
		fmt.Println("Error!", err.Error())
		return
	}


	var s = []S{}
	for _, v := range cities {
		s = append(s, S{Type: "c", Id: strconv.Itoa(v.City_id), Value: (v.City + ", " + v.Country), Count: strconv.Itoa(v.Count)})
	}
	for _, v := range hotels {
		s = append(s, S{Type: "h", Id: strconv.Itoa(v.Hotel_id), Value: (v.Hotel_name + ", " + v.City + ", " + v.Country), Rating_average: strconv.FormatFloat(v.Rating_average, 'f', -1, 64)})
	}


	var json, err2 = json.Marshal(s)

	if err2 != nil {
		http.Error(w, err2.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(json)
	return
	// }

	http.Error(w, "", http.StatusBadRequest)
}