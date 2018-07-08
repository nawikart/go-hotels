package main

import (
	"github.com/go-redis/redis"
	"strings"
	"./custom"
	"net/http"
	"encoding/json"
	"fmt"
)

type S struct {
	Tipe	string
	Id      string
	Value   string
	Count   string
	Rating_average   string
}

func autocompleteRedis(w http.ResponseWriter, r *http.Request) {

	keyword := r.URL.Query()["key"][0]
	// fmt.Println(keyword)

	pattern := "*"+ custom.String2InCaseSensitivePattern(keyword) +"*"
	// fmt.Println(pattern)

	var s = []S{}
	
	res := zScan("cities", 0, pattern, 30000)
	var vs = []string{}
	i:=0
	for _, v := range res {
		vs = strings.Split(v, ":")
		if len(vs) >= 3 {
			// fmt.Println(vs[1])
			s = append(s, S{Tipe: "c", Id: vs[0], Value: vs[1], Count: vs[2]})
			i++
			if i > 5 {
				break
			}
		}
	}

	res = zScan("hotels", 0, pattern, 1000000)
	i=0
	for _, v := range res {
		vs = strings.Split(v, ":")
		if len(vs) >= 3 {
			fmt.Println(vs)
			s = append(s, S{Tipe: "h", Id: vs[0], Value: vs[1], Rating_average: vs[2]})
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
	w.Write(js)
}


func zScan(k string, cs int, p string, count int64)  []string{
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	res, _, _ := client.ZScan(k, 0, p, count).Result()
	return res
}


func main() {

    http.HandleFunc("/autocomple-redis", autocompleteRedis)


    fmt.Println("starting web server at http://localhost:8081/")
    http.ListenAndServe(":8081", nil)
}