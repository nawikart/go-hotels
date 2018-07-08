package apiControllers

import (
    // "fmt"
    "net/http"
    "bytes"
    "io/ioutil"
)

func agoda(data string) []byte{

    // fmt.Println(data)
    
    url := "http://affiliateapi7643.agoda.com/affiliateservice/lt_v1"

    var jsonStr = []byte(data)
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
    req.Header.Set("Authorization", "1730363:e9cf134d-0b06-46e8-8957-3d4d15042194")
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    // fmt.Println("response Status:", resp.Status)
    // fmt.Println("response Headers:", resp.Header)
    body, _ := ioutil.ReadAll(resp.Body)
    

    return []byte(body)
}