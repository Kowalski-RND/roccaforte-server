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
	log.Println("Roccaforte API up on port 9090. Quit with ^C!")
	log.Fatal(http.ListenAndServe(":9090", r))
}
