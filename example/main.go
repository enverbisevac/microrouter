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

	root := microrouter.NewGroup()
	// method first and then path
	root.Add("/enver", index, "GET", "POST")

	// admin group
	adminGroup := microrouter.NewGroup()
	adminGroup.Add("/dashboard", dashboard, "GET")

	// profile group
	profile := microrouter.NewGroup()
	profile.Add("/user", user, "GET")

	// assign profile to admin
	adminGroup.AddGroup("/profile", profile)
	// assign admin to root
	root.AddGroup("/admin", adminGroup)

	router.SetRootGroup(root)
	log.Fatal(http.ListenAndServe(":8080", router))
}

func index(res http.ResponseWriter, req *http.Request) {
	_, _ = fmt.Fprint(res, "<h1>Hello from index handler</h1>")
}

func dashboard(res http.ResponseWriter, req *http.Request) {
	_, _ = fmt.Fprint(res, "<h1>Hello from dashboard</h1>")
}

func user(res http.ResponseWriter, req *http.Request) {
	_, _ = fmt.Fprint(res, "<h1>Hello from profile</h1>")
}
