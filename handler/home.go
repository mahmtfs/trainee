package handler

import (
	"html/template"
	"net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request){
	t, _ := template.ParseFiles("templates/index.html")
	err := t.Execute(w, nil)
	HandleError(err)
}
