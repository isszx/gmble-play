package search

import (
	"database/sql"
	"fmt"
	"github.com/iotku/mumzic/config"
	_ "github.com/mattn/go-sqlite3"
	"os"
)

var MaxDBID int

func init() {
	if _, err := os.Stat(config.Songdb); os.IsNotExist(err) {
		MaxDBID = 0
	} else {
		// Number of rows (not to exceed) in sqlite database
		MaxDBID = getMaxID(config.Songdb)
	}
}

// Aggressively fail on error
func checkErrPanic(err error) {
	if err != nil {
		panic(err)
	}
}

// Query SQLite database to count maximum amount of rows, as to not point to non existent ID
// TODO: Perhaps catch error instead?
func getMaxID(database string) int {
	db, err := sql.Open("sqlite3", database)
	defer db.Close()
	checkErrPanic(err)
	var count int
	err = db.QueryRow("select max(ROWID) from music;").Scan(&count)
	checkErrPanic(err)
	return count
}

// Query SQLite database to get filepath related to ID
func GetTrackById(trackID int) (filepath, humanout string) {
	if trackID > MaxDBID {
		return "", ""
	}
	db, err := sql.Open("sqlite3", config.Songdb)
	checkErrPanic(err)
	defer db.Close()
	var path, artist, title, album string
	err = db.QueryRow("select path,artist,title,album from MUSIC where ROWID = ?", trackID).Scan(&path, &artist, &title, &album)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return "", ""
		}
	}
	checkErrPanic(err)

	humanout = artist + " - " + title
	return path, humanout
}

func SearchALL(Query string) []string {
	Query = fmt.Sprintf("%%%s%%", Query)
	rows := makeDbQuery(config.Songdb, "SELECT ROWID, * FROM music where (artist || \" \" || title)  LIKE ? LIMIT 25", Query)
	defer rows.Close()

	var rowID int
	var artist, album, title, path string
	var output []string
	for rows.Next() {
		err := rows.Scan(&rowID, &artist, &album, &title, &path)
		checkErrPanic(err)
		output = append(output, fmt.Sprintf("#%d | %s - %s (%s)", rowID, artist, title, album))
	}

	return output
}

func ShowFullList() []string {
	rows := makeDbQuery(config.Songdb, "SELECT ROWID, * FROM music LIMIT 25")
	defer rows.Close()

	var rowID int
	var artist, album, title, path string
	var output []string
	for rows.Next() {
		err := rows.Scan(&rowID, &artist, &album, &title, &path)
		checkErrPanic(err)
		output = append(output, fmt.Sprintf("#%d | %s - %s (%s)", rowID, artist, title, album))
	}

	return output
}

// Helper Functions
func makeDbQuery(songdb, query string, args ...interface{}) *sql.Rows {
	db, err := sql.Open("sqlite3", songdb)
	checkErrPanic(err)
	defer db.Close()
	rows, err := db.Query(query, args...)
	checkErrPanic(err)

	// Don't forget to close in function where called.
	return rows
}
