# router
This package is a router for Golang applications serving HTTP.
The router makes you follow strict rules.

## Usage

```go
package main

import (
    "fmt"
    "net/http"
    "github.com/shumichenko/router"
    "log"
)

func main() {
    appRouter := router.NewRouter()
    routesList := []router.Route{
        router.NewRoute("/users", http.MethodGet, GetUserList),
        router.NewRoute("/users/:id", http.MethodGet, GetUser),
    }
    appRouter.AddRoutes(routesList)

    server := &http.Server{
      Handler:      appRouter,
      Addr:         "127.0.0.1:8080",
      WriteTimeout: 20 * time.Second,
      ReadTimeout:  20 * time.Second,
    }
    log.Fatal(server.ListenAndServe())
}

func GetUserList(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "Igor, Dmitry, Alexander\n")
}

func GetUser(w http.ResponseWriter, r *http.Request) {
    params := router.GetParamsFromContext(r.Context())
    fmt.Fprintf(w, "User with id: %s!\n", params.GetByName("id"))
}
```

### Fetching path parameters
Path parameters are stored in request context. 
To fetch them you should pass request context to the function `router.GetParamsFromContext(r.Context)`.
```go
// Route: /news/:title
// Requested URL: /news/hello-this-is-test-2020-03
params := router.GetParamsFromContext(r.Context())
title := params.GetByName("title") // hello-this-is-test-2020-03
```

### Routes declaration rules
- request method should NOT be empty
- path should start with "/" symbol (in case of path requested like this,
  starting slash will be automatically added to requested path)
- path should NOT be empty
- path should NOT contain trailing "/" symbol (in case of path requested like this,
trailing slash will be automatically removed from requested path)
- router is not case-sensitive, declared path will be automatically formatted to lower case
- you can't declare intersecting routes with same request method. Requested path can match only one route.  
<br />
**Example of intersecting routes:**
    ```
    GET  /users            ###
    GET  /users            ### intersecting
    GET  /users/:id
    GET  /users/:id/comments
    GET  /users/:id/statistics
    POST /users
    ```
    ```
    GET  /users
    GET  /users/:id         ###
    GET  /users/statistics  ### intersecting
    GET  /users/:id/comments
    GET  /users/:id/statistics
    POST /users
    ```
    ```
    GET  /users
    GET  /users/:id         ###
    GET  /users/:dynamic    ### intersecting
    GET  /users/:id/comments
    GET  /users/:id/statistics
    POST /users
    ```
    ```
    GET  /users
    GET  /users/:id         
    GET  /users/:id/:type     ###
    GET  /users/:id/comments  ### intersecting
    GET  /users/example/types ### intersecting
    POST /users
    ```