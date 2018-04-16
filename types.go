package main

import "database/sql"


// Types that implement the rowType interface can be used in returnQueryResults
// to get the results of a postgres query in JSON.
// These types should contain fields corresponding to the attributes of each
// row in the expected query result.
type rowType interface {
	
	// readFrom takes a sql Rows object and scans the current Row into this struct.
	readFrom(rows *sql.Rows) (rowType, error)

	// getQuery returns the query corresponding to this type.
	getQuery() (string)
}


// researcherType shows information on all of the researchers in the database.
type researcherType struct {
    Id          int     `json:"id"`
    FirstName   string  `json:"firstName"`
    LastName    string  `json:"lastName"`
    Email       string  `json:"email"`
}
func (row researcherType) readFrom(rows *sql.Rows) (rowType, error) {
	err := rows.Scan(&row.Id, &row.FirstName, &row.LastName, &row.Email)
	return row, err
}
func (row researcherType) getQuery() string {
	return "SELECT * FROM researcher;"
}


// projectType shows information on all of the projects in the database.
type projectType struct {
    Id          int     `json:"id"`
    Name        string  `json:"Name"`
    Date        string  `json:"date"`
}
func (row projectType) readFrom(rows *sql.Rows) (rowType, error) {
	err := rows.Scan(&row.Id, &row.Name, &row.Date)
	return row, err
}
func (row projectType) getQuery() string {
	return "SELECT * FROM project;"
}


// authorshipType shows the name of each project with the names of all of its authors.
type authorshipType struct {
    ProjectName string  `json:"projectName"`
    AuthorName  string  `json:"authorName"`
}
func (row authorshipType) readFrom(rows *sql.Rows) (rowType, error) {
	err := rows.Scan(&row.ProjectName, &row.AuthorName)
	return row, err
}
func (row authorshipType) getQuery() string {
	return `
SELECT project.name AS project, researcher.first_name || researcher.last_name AS author
FROM authorship
INNER JOIN project ON project_id = project.id
INNER JOIN researcher ON researcher_id = researcher.id
;`
}
