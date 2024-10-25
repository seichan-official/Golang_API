package main

import (
    "log"
    "net/http"
)

func main() {
    port := ":8080"

    fs := http.FileServer(http.Dir("./docs/dist"))
    http.Handle("/", fs)

    // サーバーの起動
    log.Printf("Server started at http://localhost%s/", port)
    err := http.ListenAndServe(port, nil)
    if err != nil {
        log.Fatal(err)
    }
}
