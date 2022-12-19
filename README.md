# router
This package is a router for Golang applications serving HTTP.
The router is not designed to have many settings and makes you follow strict rules.

### Routes declaration rules
- request method should NOT be empty
- path should start with "/" symbol (in case of path requested like this,
  starting slash will be automatically added to requested path)
- path should NOT be empty
- path should NOT contain trailing "/" symbol (in case of path requested like this,
trailing slash will be automatically removed from requested path)
- router is not case-sensitive, declared path will be automatically formatted to lower case
- you can't declare intersecting routes. Requested path can match only one route.  
**Example:**
    ```
    GET  /news            ###
    GET  /news            ### intersecting
    GET  /news/:id
    GET  /news/:id/comments
    GET  /news/:id/statistics
    POST /news
    ```
    ```
    GET  /news
    GET  /news/:id         ###
    GET  /news/statistics  ### intersecting
    GET  /news/:id/comments
    GET  /news/:id/statistics
    POST /news
    ```
    ```
    GET  /news
    GET  /news/:id         ###
    GET  /news/:dynamic    ### intersecting
    GET  /news/:id/comments
    GET  /news/:id/statistics
    POST /news
    ```
    ```
    GET  /news
    GET  /news/:id         
    GET  /news/:id/:type     ###
    GET  /news/:id/comments  ### intersecting
    GET  /news/example/types ### intersecting
    POST /news
    ```