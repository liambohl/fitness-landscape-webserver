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


type researcherType struct {
    Id          int     `json:"id"`
    FirstName   string  `json:"firstName"`
    LastName    string  `json:"lastName"`
    Email       string  `json:"email"`
}

func main() {
    fmt.Println("Starting")
    // Parse the command line arguments
    flag.StringVar(&configpath, "config", "/etc/app/dbconfig.yaml", "Override the default config file path")
    flag.Parse()

    // Assign the config values to the global config struct
    config = GetPostgresConfig()
    connection = DatabaseString(config)

    // Start the web server
    http.HandleFunc("/researcher", handleResearcher)
    http.ListenAndServe(":5000", nil)
}


func handleResearcher(resp http.ResponseWriter, req *http.Request) {
    var researchers []researcherType

    rows := QueryDatabase("SELECT * FROM researcher;")
    defer rows.Close()
    for rows.Next() {
        var researcher researcherType
        err := rows.Scan(&researcher.Id, &researcher.FirstName,
            &researcher.LastName, &researcher.Email)
        if err != nil {
            fmt.Println("Problem iterating results")
        }
        researchers = append(researchers, researcher)
    }

    json, err := json.Marshal(&researchers)
    if err != nil {
        fmt.Println(err.Error())
    }
    resp.Write(json)
}


// Function handles database queries
// Returns false if bad query
func QueryDatabase(query string) *sql.Rows {
    fmt.Println("Opening db connection")
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
