[![Actions Status](https://github.com/enverbisevac/microrouter/workflows/Go/badge.svg)](https://github.com/enverbisevac/microrouter/actions)
[![codecov](https://codecov.io/gh/enverbisevac/microrouter/branch/master/graph/badge.svg)](https://codecov.io/gh/enverbisevac/microrouter)
[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)
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
- [ ] Name for url path
- [ ] Group urls
- [ ] Static files
- [x] Default handler for some standard content types
- [x] Example 
