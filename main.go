package main

import (
	"fmt"
	"log"
	"net/http"
	never "never/HTML"
)

func main() {
	http.HandleFunc("/", never.HandleRequest)
	http.HandleFunc("/info.html", never.HandleRequest2)
	http.Handle("/templates/", http.StripPrefix("/templates/", http.FileServer(http.Dir("templates/"))))
	http.Handle("/HTML/", http.StripPrefix("/HTML/", http.FileServer(http.Dir("HTML/"))))
	log.Fatal(http.ListenAndServe(":8080", nil))
	fmt.Println("go to -->  http://localhost:8080/ ")

}
