package main
import (
    "io/ioutil"
    "fmt"
    "net/http"
    "encoding/json"
    "flag"
    "database/sql"
    _   "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres"
)

var configpath string       // Global var used in config.go -- path of postgres acess config
var config PostgresConfig   // Global config struct from config.go in main package
var db *sql.DB               // Database connection


func main() {
    fmt.Println("Starting")

    // Parse the command line arguments
    flag.StringVar(&configpath, "config", "/etc/app/dbconfig.yaml", "Override the default config file path")
    flag.Parse()

    // Configure and open database connection
    config = GetPostgresConfig()
    fmt.Println("Opening db connection...")
    db, err := sql.Open("cloudsqlpostgres", DatabaseString(config))
    defer db.Close()
    if err != nil {
        fmt.Println("Failed to connect to db")
        fmt.Println(err.Error())
    }
    fmt.Println("Connected to the database")

    // Handle each request
    http.HandleFunc("/hello", handleHello)
    http.HandleFunc("/researcher", func (resp http.ResponseWriter, req *http.Request) {
        returnQueryResults(resp, researcherType{}, db)
    })
    http.HandleFunc("/project", func (resp http.ResponseWriter, req *http.Request) {
        returnQueryResults(resp, projectType{}, db)
    })
    http.HandleFunc("/authorship", func (resp http.ResponseWriter, req *http.Request) {
        returnQueryResults(resp, authorshipType{}, db)
    })
    http.HandleFunc("/post/researcher", func (resp http.ResponseWriter, req *http.Request) {
        handlePostResearcher(resp, req, db)
    })

    // Start the web server
    http.ListenAndServe(":5000", nil)
}


// handleHello is a simple handler for debugging.
// It writes to both the http response and stdout.
func handleHello(resp http.ResponseWriter, req *http.Request) {
    fmt.Println("Hello")
    resp.Write([]byte("Hello, World!"))
}


// returnQueryResults taks an entry which holds the query and the expected
// format of each row in the result.
// It sends the query to the database, marshals the response into JSON, and
// sends the resultant JSON as an HTTP response.
func returnQueryResults(resp http.ResponseWriter, entry rowType, db *sql.DB) {
    query := entry.getQuery()
    rows := QueryDatabase(query, db)
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


// QueryDatabase opens a connection to the database and returns the result
// of a given query.
func QueryDatabase(query string, db *sql.DB) *sql.Rows {
    rows, err := db.Query(query)
    if err != nil {
        fmt.Println("Failed to query db")
        fmt.Println(err.Error())
    }
    fmt.Println("Query returned")

    return rows
}


func handlePostResearcher(resp http.ResponseWriter, req *http.Request, db *sql.DB) {
    // Get POST data
    body, err := ioutil.ReadAll(req.Body)
    if err != nil {
        fmt.Println(err.Error())
    }

    var entries []researcherType
    if err := json.Unmarshal(body, &entries); err != nil {
        fmt.Println(err.Error())
    }

    for _, entry := range(entries) {
        query := "INSERT INTO researcher(first_name, last_name, email) VALUES ($1, $2, $3);"
        _, err := db.Query(query, entry.FirstName, entry.LastName, entry.Email)
        if err != nil {
            fmt.Println("Failed to query db")
            fmt.Println(err.Error())
        }
        fmt.Println("Insert executed")
    }

    resp.Write([]byte("Data successfully written\n"))
}
