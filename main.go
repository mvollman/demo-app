package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "github.com/go-redis/redis/v7"
    "github.com/elastic/go-elasticsearch/v8"

    "database/sql"
    _ "github.com/go-sql-driver/mysql"
)

func handler(w http.ResponseWriter, r *http.Request) {

    fmt.Sprintf("New Connection")
    fmt.Fprintf(w, "<html>")

    // elasticsearch
    es, err := elasticsearch.NewDefaultClient()
    if err != nil {
        log.Fatalf("Error creating the client: %s", err)
    }
    res, err := es.Info()
    if err != nil {
      log.Fatalf("Error getting response: %s", err)
    }
    if res.IsError() {
      log.Fatalf("Error: %s", res.String())
    }
    if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
      log.Fatalf("Error parsing the response body: %s", err)
    }
    fmt.Sprintf("Elasticsearch Version: %s", elasticsearch.Version)
    fmt.Fprintf(w, "<p>Elasticsearch Version: %s</p>", elasticsearch.Version)

    // redis
    client := redis.NewClient(&redis.Options{
        Addr:     os.Getenv("REDIS_URL"),
        Password: "", // no password set
        DB:       0,  // use default DB
    })

    pong, err := client.Ping().Result()
    if err != nil {
        log.Fatalf("Error response from redis %s", err)
    }
    fmt.Sprintf("Redis pong?: %s", pong)
    fmt.Fprintf(w, "<p>Redis pong?: %s</p>", pong)


    // mysql
    db, err := sql.Open("mysql", os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASS") + "@tcp(" + os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT") + ")/" + os.Getenv("DB_NAME"))

    if err != nil {
        panic(err.Error())
    }

    defer db.Close()

    results, err := db.Query("show tables;")
    if err != nil {
        panic(err.Error())
    }

    var table string

    for results.Next() {
        results.Scan(&table)
        fmt.Sprintf("MYSQL query results: %s", table)
        fmt.Fprintf(w, "<p>MYSQL query results: %s</p>", table)
    }

    fmt.Sprintf("Connection Closed")
    fmt.Fprintf(w, "</html>")
}

func main() {
    fmt.Sprintf("Starting Demo APP")
    http.HandleFunc("/", handler)
    log.Fatal(http.ListenAndServe(":8080", nil))
}
