package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
)

func main() {
	http.HandleFunc("/", SubmitHandler)

	// log.Fatal(http.ListenAndServe("localhost:8000", nil))

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

/*
func RootHandler(w http.ResponseWriter, r *http.Request) {
	// r.ParseForm()

	fmt.Println(r.Form)
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	fmt.Fprintf(w, "Hello astaxie!")
}
*/
