package main
import (
    "fmt"
    "net/http"
    "encoding/json"
    "flag"
    "database/sql"
    //_ "github.com/lib/pq"
    _   "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres"
)

var configpath string       // Global var used in config.go -- path of postgres acess config
var config PostgresConfig   // Global config struct from config.go in main package
var connection string       // username, password, host, etc. to access postgres server


func main() {
    fmt.Println("Starting")
    // Parse the command line arguments
    flag.StringVar(&configpath, "config", "/etc/app/dbconfig.yaml", "Override the default config file path")
    flag.Parse()

    // Assign the config values to the global config struct
    config = GetPostgresConfig()
    connection = DatabaseString(config)

    // Start the web server
    http.HandleFunc("/hello", handleHello)
    http.HandleFunc("/researcher", func (resp http.ResponseWriter, req *http.Request) {
        returnQueryResults(resp, researcherType{})
    })
    http.HandleFunc("/project", func (resp http.ResponseWriter, req *http.Request) {
        returnQueryResults(resp, projectType{})
    })
    http.HandleFunc("/authorship", func (resp http.ResponseWriter, req *http.Request) {
        returnQueryResults(resp, authorshipType{})
    })
    http.ListenAndServe(":5000", nil)
}


func handleHello(resp http.ResponseWriter, req *http.Request) {
    fmt.Println("Hello")
    resp.Write([]byte("Hello, World!"))
}


func returnQueryResults(resp http.ResponseWriter, entry rowType) {
    query := entry.getQuery()
    rows := QueryDatabase(query)
    defer rows.Close()

    var entries []interface{}
    var err error
    for rows.Next() {
        entry, err = entry.readFrom(rows)
        if err != nil {
            fmt.Println("Problem iterating results")
        }
        entries = append(entries, entry)
    }

    json, err := json.Marshal(&entries)
    if err != nil {
        fmt.Println(err.Error())
    }
    resp.Write(json)
}


// Function handles database queries
// Returns false if bad query
func QueryDatabase(query string) *sql.Rows {
    fmt.Println("Opening db connection...")
    db, err := sql.Open("cloudsqlpostgres", connection)
    defer db.Close()
    if err != nil {
        fmt.Println("Failed to connect to db")
        fmt.Println(err.Error())
    }
    fmt.Println("Connected to the database")


    rows, err := db.Query(query)
    if err != nil {
        fmt.Println("Failed to query db")
        fmt.Println(err.Error())
    }
    fmt.Println("Query returned")


    return rows
}
