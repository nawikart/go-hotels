package main
 
import (
	"fmt"
	"net/http"
	"io/ioutil"
	"regexp"
	"./db/mysql"
	"strings"
)
 

func curlAgoReviews(hotel_id string, url string){

	mysql := mysql.Connect()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
	}


	res, err2 := http.DefaultClient.Do(req)
 	if err2 != nil {
		fmt.Println(err2)
	}

	defer res.Body.Close()
	
	body, _ := ioutil.ReadAll(res.Body)

	rg_newline := regexp.MustCompile("\\n+")
	rg_singlequote := regexp.MustCompile("'")
	rg_2spaces := regexp.MustCompile("\\s\\s+")
	body_ok := string(rg_2spaces.ReplaceAllString(rg_newline.ReplaceAllString(string(body), " "), " "))

	// fmt.Println(body_ok)

	//, reviewPageParams.hotelInfo = { hotelId: 1166068, name: "
	//fmt.Println("=======HOTEL INFO=======")
	rg_hotel, _ := regexp.Compile(`, reviewPageParams.hotelInfo = { hotelId: ([0-9]+), name: "`)
	res_hotel := rg_hotel.FindStringSubmatch(body_ok)
	
	
	//fmt.Println("=======RATING=======")
	rg_rating, _ := regexp.Compile(`reviewPageParams.reviews = { scoreText: "(.*)", score: "(.*)", basedOn: "(.*)", reviewsCount: (.*), recommendationScore: (.*), isShowRecommendationScore: (.*), isReviewPage: (.*), }; reviewPageParams.mapOpenTrackUrl`)
	res_rating := rg_rating.FindStringSubmatch(body_ok)
	

	//fmt.Println("=======ADDRESS=======")
	rg_address, _ := regexp.Compile(`, address: { full: "(.*)", countryName: "(.*)", areaName: "(.+)", address: "(.*)", postalCode: "(.*)", cityName: "(.*)", cityId: (.*), }, isGetQuestionEnabled`)
	res_address := rg_address.FindStringSubmatch(body_ok)
	
	
	//fmt.Println("=======REVIEWS=======")
	rg_reviews, _ := regexp.Compile(`reviews: {"commentList":{"provider":332,"comments":\[(.+)\],"haveMoreThanOneComments`)
	res_reviews := rg_reviews.FindStringSubmatch(body_ok)
	// fmt.Println(res_reviews)

	if len(res_hotel) > 0 && len(res_rating) > 0 && len(res_address) > 0 {
		// fmt.Println(res_reviews)
		var reviewsJson string
		if len(res_reviews) > 0 {
			reviewsJson = "[" + res_reviews[1] + "]"
		}
		sql_ins_into_hotels_custom := `INSERT IGNORE 
			INTO hotels_custom (
			hotel_id,
			scoreText, 
			score, 
			reviewsCount, 
			recommendationScore, 
			address_countryName, 
			address_areaName, 
			address_address, 
			address_postalCode, 
			address_cityName, 
			address_cityId, 
			address_full,
			reviews) 
			VALUES (
			'`+ res_hotel[1] +`',
			'`+ res_rating[1] +`',
			'`+ res_rating[2] +`',
			'`+ res_rating[4] +`',
			'`+ res_rating[5] +`',
			'`+ res_address[2] +`',
			'`+ res_address[3] +`',
			'`+ res_address[4] +`',
			'`+ res_address[5] +`',
			'`+ res_address[6] +`',
			'`+ res_address[7] +`',
			'`+ res_address[1] +`',
			'`+ rg_singlequote.ReplaceAllString(reviewsJson, `\'`) +`'
			)`
		
		// fmt.Println(sql_ins_into_hotels_custom)
		// fmt.Println("===========sql_ins_into_hotels_custom============")
	
		mysql.Query(sql_ins_into_hotels_custom)
		mysql.Query("UPDATE hotels_link_to_agoda SET status = '1' WHERE hotel_id = '" + res_hotel[1] + "'")
	
		fmt.Println("\nsukses grab hotel ID: ", res_hotel[1], " ::URL -->", url)
		fmt.Println("\n")
		
	}else{

		// fmt.Println("=======HOTEL INFO=======")
		// fmt.Println(res_hotel)

		// fmt.Println("=======RATING=======")
		// fmt.Println(res_rating)

		// fmt.Println("=======ADDRESS=======")
		// fmt.Println(res_address)

		fmt.Println("FAILED grab URL -->", url)
		// mysql.Query("UPDATE hotels_link_to_agoda SET status = '0' WHERE hotel_id = '" + hotel_id + "'")

		// fmt.Println("=======END=======")
		// fmt.Println("==============")
		fmt.Println("")
	}
}

func main() {
	
	fmt.Println("start...")

	mysql := mysql.Connect()

 	rows, err := mysql.Query("SELECT hotel_id, hotel_namekey, citykey, countryisocode FROM hotels_link_to_agoda WHERE status IS NULL limit 0, 5")
	if err != nil {
		fmt.Println(err)
	}
	
	var url, hotel_id, hotel_namekey, citykey, countryisocode string

	for rows.Next() {
		err = rows.Scan(&hotel_id, &hotel_namekey, &citykey, &countryisocode)
		if err != nil {
			fmt.Print(err)
		}
		url = "https://www.agoda.com/" + strings.ToLower(hotel_namekey + "/reviews/"+ citykey + "-" + countryisocode +".html")
		// fmt.Println(url)

		curlAgoReviews(hotel_id, url)
	}
}