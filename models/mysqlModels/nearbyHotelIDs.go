package mysqlModels

import (
	"../../db/mysql"
	"fmt"
)

func NearbyHotelIDs(lat string, lng string) []string {
	mysql := mysql.Connect()

	//sqlStatement := "select hotel_id, ( 3959 * acos( cos( radians(" + lat + ") ) * cos( radians( latitude ) ) * cos( radians( longitude ) - radians(" + lng + ") ) + sin( radians(" + lat + ") ) * sin( radians( latitude ) ) ) ) AS jarak from hotel_ll where latitude != '' and longitude != '' and ( 3959 * acos( cos( radians(" + lat + ") ) * cos( radians( latitude ) ) * cos( radians( longitude ) - radians(" + lng + ") ) + sin( radians(" + lat + ") ) * sin( radians( latitude ) ) ) ) < 5 ORDER BY ( 3959 * acos( cos( radians(" + lat + ") ) * cos( radians( latitude ) ) * cos( radians( longitude ) - radians(" + lng + ") ) + sin( radians(" + lat + ") ) * sin( radians( latitude ) ) ) ) ASC LIMIT 0,50"
	sqlStatement := "select hotel_id from hotels_ll where latitude != '' and longitude != '' and ( 3959 * acos( cos( radians(" + lat + ") ) * cos( radians( latitude ) ) * cos( radians( longitude ) - radians(" + lng + ") ) + sin( radians(" + lat + ") ) * sin( radians( latitude ) ) ) ) < 5 ORDER BY ( 3959 * acos( cos( radians(" + lat + ") ) * cos( radians( latitude ) ) * cos( radians( longitude ) - radians(" + lng + ") ) + sin( radians(" + lat + ") ) * sin( radians( latitude ) ) ) ) ASC LIMIT 0,24"
	rows, err := mysql.Query(sqlStatement)

	if err != nil {
		fmt.Println(err)
	}
	// defer rows.Close()

	var hotel_id string
	var hotel_ids = []string{}

	for rows.Next() {
		err2 := rows.Scan(&hotel_id)
		if err2 != nil {
			fmt.Print(err2)
		}

		hotel_ids = append(hotel_ids, hotel_id)
	}

	// rows.Close()
	return hotel_ids
}
