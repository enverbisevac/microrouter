package main

import (
	"fmt"
	"github.com/enverbisevac/microrouter"
	"log"
	"net/http"
)

func main() {
	router := microrouter.NewRouter()
	router.Use(microrouter.LoggerMiddleware())
	err := router.Add("GET /enver", index)
	if err != nil {
		panic("Error adding route to router")
	}
	log.Fatal(http.ListenAndServe(":8080", router))
}

func index(res http.ResponseWriter, req *http.Request) {
	log.Println("index handler")
	fmt.Fprint(res, "<h1>Hello from index handler</h1>")
}
