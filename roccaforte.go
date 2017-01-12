package main

import (
	"github.com/roccaforte/server/api"
	"github.com/roccaforte/server/model"
	"log"
	"net/http"
)

func main() {
	model.InitDB("user=postgres dbname=roccaforte sslmode=disable")
	r := api.New()
	log.Println("Roccaforte API up on port 8080. Quit with ^C!")
	log.Fatal(http.ListenAndServe(":8080", r))
}
