package main

import (
      "os"
      "fmt"
      "io/ioutil"
      "gopkg.in/yaml.v2"
)


type PostgesConfig struct {
      Dbhostname   string    // IP address or fully qualified hostname
      Dbport       string    // The port (usually 5432 for postgres)
      Dbuser       string    // Database username
      Dbpass       string    // User's password
      Databasename string    // Name of the database
}


// Read the config file and return a PostgresConfig struct
func GetPostgresConfig() PostgesConfig {

      // Read the raw config file into memory
      yamlFile, err := ioutil.ReadFile(configpath)
      if err != nil {
           fmt.Println(err.Error())
           os.Exit(1) // Quit with exit status 1
      }


      // Unmarshal the yaml config into a PostgesConfig
      var config PostgesConfig
      err = yaml.Unmarshal(yamlFile, &config)
      if err != nil {
            fmt.Println(err.Error())
            os.Exit(1)
      }

      return config
}


// Return a Postgres connection string
func DatabaseString(c PostgesConfig) string {
      return "postgres://" + c.Dbuser + ":" + c.Dbpass + "@" + c.Dbhostname + ":" + c.Dbport + "/" + c.Databasename
}