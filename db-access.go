package main
import (
	"fmt"
	"net/http"
	"encoding/json"
	"flag"
	"database/sql"//; _ "github.com/lib/pq";
)

var configpath string		// Global var used in config.go -- path of postgres acess config
var config PostgesConfig	// Global config struct from config.go in main package
var connection string		// username, password, host, etc. to access postgres server


type researcherType struct {
	Id			int		`json:"id"`
	FirstName	string	`json:"firstName"`
	LastName	string	`json:"lastName"`
	Email		string	`json:"email"`
}

func main() {
	// Parse the command line arguments
	flag.StringVar(&configpath, "config", "/etc/app/dbconfig.yaml", "Override the default config file path")
	flag.Parse()

	// Assign the config values to the global config struct
	config = GetPostgresConfig()
	
	// Start the web server
	http.HandleFunc("/researcher", handleResearcher)
	http.ListenAndServe(":5000", nil)
}


func handleResearcher(resp http.ResponseWriter, req *http.Request) {
	var researchers []researcherType

	rows := QueryDatabase("SELECT * FROM researchers;")
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
	db, err := sql.Open("postgres", connection)
	defer db.Close()
	if err != nil {
		fmt.Printf("Failed to connect to db")
		fmt.Println(err.Error())
	}
	fmt.Println("Connected to the database")


	rows, err := db.Query(query)
	if err != nil {
		fmt.Printf("Failed to connect to db")
		fmt.Println(err.Error())
	}


	return rows
}
