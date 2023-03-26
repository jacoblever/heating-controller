package main

import (
	"fmt"
	"net/http"

	"github.com/jacoblever/heating-controller/brain/brain"
)

var port = 8080

func main() {
	router := brain.CreateRouter(brain.DefaultConfig)

	fmt.Println(fmt.Sprintf("Listening on port %d...", port))
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), router)
	if err != nil {
		panic(err)
	}
}
