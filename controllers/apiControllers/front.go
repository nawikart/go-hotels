package apiControllers

import (
	"encoding/json"	
	"net/http"
	"strings"
	// "fmt"
	"regexp"
)

func Front(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	mod := strings.ToLower(r.URL.Query()["mod"][0])
	var js []byte
	var err error
	var where string

	// fmt.Println(mod)

    switch mod {

    	case "cityurl":

    		ty := r.URL.Query()["ty"][0]
    		id := r.URL.Query()["id"][0]

    		js, err = json.Marshal(autocompleteCityUrl(ty, id))

    	case "agoda":
    		
    		js = agoda(r.URL.Query()["data"][0])

    	case "hotel":

			if(r.URL.Query()["hotel_id"] != nil){
				where = "WHERE hotel_id = "+ r.URL.Query()["hotel_id"][0]
			}else{

				hotelnamekey := r.URL.Query()["hotelnamekey"][0]
				city_cc := r.URL.Query()["city_cc"][0]

				rg_cc, _ := regexp.Compile("([-a-z0-9]+)-([a-z][a-z]).html")
				url_cc := rg_cc.FindStringSubmatch(city_cc)
				citykey := url_cc[1]
				// countryisocode := url_cc[2]

				where = "WHERE hotel_namekey = '" + hotelnamekey + "' AND citykey = '" + citykey + "'"
			}

    		js, err = json.Marshal(hotel(where))

    	case "hotels":

    		if(r.URL.Query()["hids"] != nil){
    			hids := r.URL.Query()["hids"][0]
    			citykey := r.URL.Query()["citykey"][0]
    			js, err = json.Marshal(hotels("", hids, citykey, "", ""))
    		}else{
	    		city_cc := r.URL.Query()["city_cc"][0]

				rg_cc, _ := regexp.Compile("([-a-z0-9]+)-([a-z][a-z])-hotels((-p([0-9]+))*).html")
				url_cc := rg_cc.FindStringSubmatch(city_cc)

				// fmt.Println(url_cc)
				dataFilter := ""
				citykey := url_cc[1]
				countryisocode := url_cc[2]

				if(r.URL.Query()["dataFilter"] != nil){
					dataFilter = r.URL.Query()["dataFilter"][0]
				}

				page := "1"
				if url_cc[5] != "" {
					page = url_cc[5]
				}

				js, err = json.Marshal(hotels(dataFilter, "", citykey, countryisocode, page))
			}


    	case "cities":

    		var countrykey string
			if(r.URL.Query()["countrykey"] != nil){
				countrykey = r.URL.Query()["countrykey"][0]
			}else{
				countrykey = ""
			}       		
    		js, err = json.Marshal(cities(countrykey))

    	case "gmap":

    		city_cc := r.URL.Query()["city_cc"][0]
			rg_cc, _ := regexp.Compile("([-a-z0-9]+)-([a-z][a-z]).html")
			url_cc := rg_cc.FindStringSubmatch(city_cc)
			citykey := url_cc[1]
			countryisocode := url_cc[2]

    		js, err = json.Marshal(gmap(citykey, countryisocode))
    }	
	
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(js)
}