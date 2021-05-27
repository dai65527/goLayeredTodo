package handler

import (
	"fmt"
	"net/http"
	"strings"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "<h1>this is todo api server</h1>")
}

func getRouteParams(r *http.Request) []string {
	splited := strings.Split(r.RequestURI, "/")
	var params []string
	for i := 0; i < len(splited); i++ {
		if len(splited[i]) != 0 {
			params = append(params, splited[i])
		}
	}
	return params
}
