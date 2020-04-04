package main

import (
	"database/sql"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"

	gUS "github.com/damilarelana/goUrlShortener"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

// initialize the database connection parameters
// * these are used default values for testing purposes
var (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "brainiac"
	database = "go_test_db"
)

// // DatabaseConnectionParameter types
// type DatabaseConnectionParameter struct {
// 	Host     string
// 	Port     int
// 	UserName string
// 	Password string
// 	DBName   string
// }

// defines the error message handler
func errMsgHandler(msg string, err error) {
	if err != nil {
		fmt.Println(msg, err.Error())
		os.Exit(1)
	}
}

// define flags
var yamlFilename *string = flag.String("yaml", "", "a yaml file containing path and mapped URL, in a 'question, answer' format per record line")
var jsonFilename *string = flag.String("json", "", "a json file containing path and mapped URL, in a 'question, answer' format per record line")
var sqlDatabasePath *string = flag.String("sql", "", "an sql database path to 'question, answer' records, with `path` and mapped `URL` in table columns per record")

// fileReader()
//  * takes the pointer to the opened file
//	* dereferences the pointer to get the value i.e. *f
//  * use ioutil to read the file and return data as []bytes
//  * returns the function call as []bytes
func fileReader(f *string) []byte {
	data, err := ioutil.ReadFile(*f)
	errMsgHandler(fmt.Sprintf("Failed to read file: %s\n", *f), err)
	return data
}

// sqlFlagReader()
//  * takes the pointer to the sql database path provided by the user in the format:
//			- `<protocol>://<host>:<port>/<database>?<username>&<password>`
//	* dereferences the pointer to get the value i.e. *f
//	* uses regex to extract the following:
//			- protocol
//			- host
//			- port
//			- database i.e. dbname
//			- username
//			- password
//  * returns a `string` of database connection parameters required by sql.Open() based on the format below
// 	 		- "user:password@tcp(localhost:port)/dbname" when using "mysql"
// 			- "host=%s port=%d user=%s password=%s dbname=% sslmode=disable" when using "postgres"
func sqlFlagReader(sqlDatabasePath *string) (dbConnParams string) {
	dbConnParams = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, database)
	return dbConnParams
}

// dbQuery()
//	* takes the pointer to the sql database path
//	* uses the sqlFlagReader() to extract database connection parameters, using the sql database path
//  * reads the records from the SQL database
//  * returns the data as a []byte
func dbQuery(sqlDatabasePath *string) []gUS.PathURL {
	dbConnParams := sqlFlagReader(sqlDatabasePath)
	db := dbConnect(dbConnParams) // initiate connection to the database

	// query for multiple records from database
	var pathUrls gUS.PathURL                    // declare pathUrls struct, meant to read database data into
	var pathUrlsSlice []gUS.PathURL             // declare slice of pathUrls, i.e. of aggregations of read database data
	sqlStatement := `select * from paths`       // get all columns from all records, in table `paths`
	multipleRows, err := db.Query(sqlStatement) // execute the sql
	errMsgHandler(fmt.Sprintf("Failed to query the database"), err)

	defer multipleRows.Close() // needed in case this did not go well

	for multipleRows.Next() { // start iterating over the returned rows i.e.
		err = multipleRows.Scan(&pathUrls.Path, &pathUrls.URL)
		errMsgHandler(fmt.Sprintf("Failed to store sql data inside query the database"), err)
		pathUrlsSlice = append(pathUrlsSlice, pathUrls) // append the newly scanned user to the existing slice of Users
	}

	err = multipleRows.Err() // handle the error thrown when `multipleRows.Next() returns a false`
	errMsgHandler(fmt.Sprintf("multipleRows.Next() returned a false"), err)

	return pathUrlsSlice // return by aggregated data
}

// dbConnect()
// * connects to the database
// * adds the connection to a poll
// * pings the opened connection to ensure that it is still working
// * returns a db connection
func dbConnect(dbConnParams string) *sql.DB {
	db, err := sql.Open("postgres", dbConnParams) // this opens a connection and adds to the pool
	errMsgHandler(fmt.Sprintf("Failed to connect to the database"), err)
	defer db.Close()

	// connect to the database
	err = db.Ping() // this validates that the opened connection "db" is actually working
	errMsgHandler(fmt.Sprintf("The database connection is no longer open"), err)

	fmt.Println("Successfully connected to the database")
	return db
}

// multiFlagTester()
//  * checks to see if the user is trying to use multiple file formats (yaml, json, sql) at the same time
//  * check to set number of flags `n` that have been set
//	* throws an error message, if `n` > 1
//  * returns a true boolean, if `n` > 1
func multiFlagTester() bool {
	numFlag := flag.NFlag()
	if numFlag > 1 {
		return true
	}
	return false
}

// yamlFlagHandler()
func yamlFlagHandler(yamlFilename *string, mapHandler http.HandlerFunc) http.HandlerFunc {
	yamlHandler, err := gUS.YAMLHandler(fileReader(yamlFilename), mapHandler)
	errMsgHandler(fmt.Sprintf("Failed to parse the YAML"), err)
	return yamlHandler
}

// jsonFlagHandler()
func jsonFlagHandler(jsonFilename *string, mapHandler http.HandlerFunc) http.HandlerFunc {
	jsonHandler, err := gUS.JSONHandler(fileReader(jsonFilename), mapHandler)
	errMsgHandler(fmt.Sprintf("Failed to parse the JSON"), err)
	return jsonHandler
}

// sqlFlagHandler()
func sqlFlagHandler(sqlDatabasePath *string, mapHandler http.HandlerFunc) http.HandlerFunc {
	sqlHandler, err := gUS.SQLHandler(dbQuery(sqlDatabasePath), mapHandler)
	errMsgHandler(fmt.Sprintf("Failed to read from database"), err)
	return sqlHandler
}

// formatHandler()
// * validates to avoid multiple flags being used simultaneously
// * checks for which flag is being used
// * leverages the appropriate flag handler
// * defaults to using yaml flag when no flag is chosen
func selectFlagHandler(mapHandler http.HandlerFunc) http.HandlerFunc {
	if multiFlagTester() { // check if multiple flags are being used i.e. multiFlagTester returns a true or false boolean
		errMsgHandler(fmt.Sprintf("Cannot use multiple flags at once. Please choose only yaml or json or sql \n"), nil)
	}

	if flag.NFlag() != 0 && !reflect.DeepEqual(jsonFilename, "") { // if json filename is used THEN yaml/sql flag would be empty string
		jsonHandler := jsonFlagHandler(jsonFilename, mapHandler)
		fmt.Printf("Now using the JSON flag with the file: %s\n", *jsonFilename)
		return jsonHandler
	}

	if flag.NFlag() != 0 && !reflect.DeepEqual(yamlFilename, "") { // if yaml filename is used THEN json/sql flags would be empty string
		yamlHandler := yamlFlagHandler(yamlFilename, mapHandler)
		return yamlHandler
	}

	if flag.NFlag() != 0 && !reflect.DeepEqual(sqlDatabasePath, "") { // if sql flag is used THEN yaml/json flags would be empty string
		sqlHandler := sqlFlagHandler(sqlDatabasePath, mapHandler)
		return sqlHandler
	}

	// defaults to yamlHandler when no flags are selected i.e. 'flag.NFlag() returns 0'
	*yamlFilename = "pathsData.yaml"
	fmt.Printf("No file flags was set. Defaulting to file: %s\n", *yamlFilename)
	yamlHandler := yamlFlagHandler(yamlFilename, mapHandler)
	return yamlHandler
}

// urlShortenerHomepage handler
func urlShortenerHomePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		custom404PageHandler(w, r, http.StatusNotFound)
		return
	}
	dataHomePage := "Url Shortener: homepage"
	io.WriteString(w, dataHomePage)
}

// custom404PageHandler defines custom 404 page
func custom404PageHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.Header().Set("Content-Type", "text/html") // set the content header type
	w.WriteHeader(status)                       // this automatically generates a 404 status code
	if reflect.DeepEqual(status, http.StatusNotFound) {
		data404Page := "This page does not exist ... 404!" // custom error message content
		io.WriteString(w, data404Page)
	}
}

// defaultMux defines the router Mux that:
//   * initializes a new Mux
//   * maps routes to handlers
func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", urlShortenerHomePage)
	return mux
}

// define main function that:
//   * uses defaultMux()
//   * uses mapHandler from `goURlShortner` package
//   * uses yamlHandler from `goURlShortner` package
//   * uses jsonHandler from `goURlShortner` package
func main() {
	// initialize all flags
	flag.Parse()

	// create an instance of defaultMux()
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := gUS.MapHandler(pathsToUrls, mux)
	fmt.Println("Starting the server on :8080")
	log.Fatal(errors.Wrap(http.ListenAndServe(":8080", selectFlagHandler(mapHandler)), "Failed to start WebServer"))
}
