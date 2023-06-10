package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
)

func main() {
	http.HandleFunc("/", SubmitHandler)

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func SubmitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("templates/submit.gtpl")
		t.Execute(w, nil)
	} else {
		r.ParseForm()

		fmt.Println("Request: ", r.Form.Get("data"))
		t, _ := template.ParseFiles("templates/submit.gtpl")
		outputText := map[string]string{"outputText": grep(r.Form.Get("data"))}

		t.Execute(w, outputText)

		fmt.Println("-----------------------------------------------------------------------------")
	}
}
