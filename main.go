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

type projectType struct {
    Id          int     `json:"id"`
    Name        string  `json:"Name"`
    Date        string  `json:"date"`
}

type authorshipType struct {
    AuthorName  string  `json:"authorName"`
    ProjectName string  `json:"projectName"`
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
    http.HandleFunc("/hello", handleHello)
    http.HandleFunc("/researcher", handleResearcher)
    http.HandleFunc("/project", handleProject)
    http.HandleFunc("/authorship", handleAuthorship)
    http.ListenAndServe(":5000", nil)
}


func handleResearcher(resp http.ResponseWriter, req *http.Request) {
    fmt.Println("Handling request...")
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


func handleProject(resp http.ResponseWriter, req *http.Request) {
    fmt.Println("Handling request...")
    var projects []projectType

    rows := QueryDatabase("SELECT * FROM project;")
    defer rows.Close()
    for rows.Next() {
        var project projectType
        err := rows.Scan(&project.Id, &project.Name, &project.Date)
        if err != nil {
            fmt.Println("Problem iterating results")
        }
        projects = append(projects, project)
    }

    json, err := json.Marshal(&projects)
    if err != nil {
        fmt.Println(err.Error())
    }
    resp.Write(json)
}


func handleAuthorship(resp http.ResponseWriter, req *http.Request) {
    fmt.Println("Handling request...")
    var authorships []authorshipType

    rows := QueryDatabase(`
SELECT project.name AS project, researcher.first_name || researcher.last_name AS auth
or
FROM authorship
INNER JOIN project ON project_id = project.id
INNER JOIN researcher ON researcher_id = researcher.id
;`)
    defer rows.Close()
    for rows.Next() {
        var authorship authorshipType
        err := rows.Scan(&authorship.AuthorName, &authorship.ProjectName)
        if err != nil {
            fmt.Println("Problem iterating results")
        }
        authorships = append(authorships, authorship)
    }

    json, err := json.Marshal(&authorships)
    if err != nil {
        fmt.Println(err.Error())
    }
    resp.Write(json)
}


func handleHello(resp http.ResponseWriter, req *http.Request) {
    fmt.Println("Hello")
    resp.Write([]byte("Hello, World!"))
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
