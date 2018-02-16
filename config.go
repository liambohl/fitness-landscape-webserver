package main

import (
    "os"
    "fmt"
    "io/ioutil"
    "gopkg.in/yaml.v2"
)


// Struct for storing acces information for the postgres database
type PostgresConfig struct {
    Dbhostname   string     // IP address or fully qualified hostname
    Dbport       string     // The port (usually 5432 for postgres)
    Dbuser       string     // Database username
    Dbpass       string     // User's password
    Databasename string     // Name of the database
}


// GetPostgresConfig reads the config file at configpath, parses it
// as yaml, and returns a PostgresConfig struct
func GetPostgresConfig() PostgresConfig {

    // Read the raw config file into memory
    yamlFile, err := ioutil.ReadFile(configpath)
    if err != nil {
        fmt.Println(err.Error())
        os.Exit(1) // Quit with exit status 1
    }

    // Unmarshal the yaml config into a PostgresConfig
    var config PostgresConfig
    err = yaml.Unmarshal(yamlFile, &config)
    if err != nil {
        fmt.Println(err.Error())
        os.Exit(1)
    }

    return config
}


// DatabaseString takes a struct with all of the database connection information
// and builds it into a Postgres connection string
func DatabaseString(c PostgresConfig) string {
    return fmt.Sprintf("host=%s dbname=%s user=%s password=%s sslmode=disable",
        c.Dbhostname, c.Databasename, c.Dbuser, c.Dbpass)
}
