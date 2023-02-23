package main

import (
	"fmt"
	"log"
	"net/http"

	"googleauth/service"
)

func main() {

	http.HandleFunc("/", service.Main)
	http.HandleFunc("/login", service.Login)
	http.HandleFunc("/callback", service.CallBack)

	fmt.Println("Started running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
