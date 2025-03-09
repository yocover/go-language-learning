package main

import (
	"log"

	"github.com/yocover/global-toolkit/net/resty"
)

func main() {

	resp, err := resty.Get("http://localhost:8080/")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(resp))
}
