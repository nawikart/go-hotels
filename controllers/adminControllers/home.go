package adminControllers

import (
	"fmt"
	"html/template"
	"net/http"
)

func Home(w http.ResponseWriter, r *http.Request) {
	var data = map[string]string{
		"Name":    "john wick",
		"Message": "have a nice day",
	}

	var t, err = template.ParseFiles("views/admin/index.html")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	t.Execute(w, data)
}
