# MicroRouter

Simple router using regex and cache regex compiled result for fast matching.
Router contains handler setters for 404, 405 and 500 errors based on Content-Type.
Middlewares are very simple to use. Standard middlewares are included in library but
you can make one very easily.

### Instalation
`go get github.com/enverbisevac/microrouter`

### Usage
    router := microrouter.NewRouter()
 	router.Use(microrouter.LoggerMiddleware())
 	err := router.Add("/hello", index, "GET", "POST")
 	if err != nil {
 		panic("Error adding route to router")
 	}
 	log.Fatal(http.ListenAndServe(":8080", router))

### TODO
- Name for url path
- Group urls
- Static files
- Default handler for some standard content types 
