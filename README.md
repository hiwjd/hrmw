# hrmw
A middleware method for [httprouter](#https://github.com/julienschmidt/httprouter) like nodejs express

## usage
```go
package main

import (
    "fmt"
    "github.com/julienschmidt/httprouter"
    "github.com/swordwinter/hrmw"
    "net/http"
    "log"
)

func wel(w http.ResponseWriter, r *http.Request, ps httprouter.Params, m *hrmw.Middleware) {
    fmt.Fprint(w, "Wel")
    m.Next(w, r, ps)
}

func come(w http.ResponseWriter, r *http.Request, _ httprouter.Params, m *hrmw.Middleware) {
    fmt.Fprintf(w, "come!\n")
}

func main() {
    router := httprouter.New()
    router.GET("/", hrmw.Use(wel, come))

    log.Fatal(http.ListenAndServe(":8080", router))
}
```
