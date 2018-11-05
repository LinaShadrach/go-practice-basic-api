package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"os"
	"net/url"

	_ "github.com/lib/pq"
)

var dbconn *sql.DB

type configuration struct {
	PORT  int `json:"port"`
	DBCONNSTR string `json:"db_conn_str"`
}

type viewer struct {
	ID    int64  `db:"id"`
	Name  string `db:"name"`
	Count int64  `db:"count"`
}

// main is intended to be called with ONE command line arg, the configuration file
func main() {
	if len(os.Args) != 2 {
		log.Fatal("Must pass the configuration file as a command line arg")
	}
	config := loadConfig(os.Args[1])
	log.Print("Connecting to DB")
	db, err := sql.Open("postgres", config.DBCONNSTR)
	if err != nil {
		log.Fatalf("Failed to connect to the DB: %s", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping the DB: %s", err)
	}
	log.Print("Connected to the database")
	dbconn = db

	log.Printf("Starting to serve traffic on port %d", config.PORT)
	http.HandleFunc("/count", handle)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", config.PORT), nil); err != nil {
		log.Fatalf("Error while serving traffic: %s", err)
	}
	log.Print("Server shutdown")
}

func handle(w http.ResponseWriter, r *http.Request) {
	log.Print("Starting request")
	defer log.Print("Finished request")
	handleCountRequest(w, r)
}

func handleCountRequest(w http.ResponseWriter, r *http.Request) {
	var err error
	var rows *sql.Rows
	key, value, paramsAreValid := processParams(w, r.URL.Query());
	if !paramsAreValid {
		return
	}
	rows, err = queryDB(key, value)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, err, "Failed to query the DB: ")
		return
	}
	defer rows.Close()
	total, errorInProcessCountQueryResults := processCountQueryResults(rows, w);
	if errorInProcessCountQueryResults {
		return
	}
	if err := rows.Err(); err != nil {
		writeResponse(w, http.StatusInternalServerError, err, "Error while interacting with the DB:")
		return
	}
	fmt.Fprintf(w, "%d", total)
}

func processParams(w http.ResponseWriter, params url.Values) (string, string, bool){
	allParamsValid := true
	key, value := "", ""
	if len(params) > 1 {
		writeResponse(w, http.StatusBadRequest, fmt.Errorf("number of parameters: %d", len(params)), "Invalid number of parameters")
		allParamsValid = false
	} else {
		if len(params) > 0 {
			var values []string
			var isValid bool
			key, values = getParam(params, w)
			value, isValid = isValidParam(key, values, w)
			allParamsValid = isValid && allParamsValid
		}
	}
	return key, value, allParamsValid
}

func queryDB(key string, value string) (*sql.Rows, error) {
	var err error
	var rows *sql.Rows
	query := "select count from viewers"
	if key == "id" {
		var id int
		id, err = strconv.Atoi(value)
		rows, err = dbconn.Query(query+" where id = $1", id)
	} else if key == "name" {
		rows, err = dbconn.Query(query+" where name = $1", value)
	} else {
		rows, err = dbconn.Query(query)
	}
	return rows, err
}

func processCountQueryResults(rows *sql.Rows, w http.ResponseWriter) (int, bool) {
	var count, total int
	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			writeResponse(w, http.StatusInternalServerError, err, "Error while iterating over rows: %s")
			return total, true
		}
		total += count
	}
	return total, false
}

func getParam(params url.Values, w http.ResponseWriter) (string, []string) {
	valKeys := reflect.ValueOf(params).MapKeys()
	keys := make([]string, len(params))
	for i := 0; i < len(params); i++ {
		keys[i] = valKeys[i].String()
	}
	return keys[0], params[keys[0]]
}

func isValidParam(key string, values []string, w http.ResponseWriter) (string, bool) {
	if len(values) != 1 {
		writeResponse(w, http.StatusBadRequest, fmt.Errorf("number of values: %d", len(values)), fmt.Sprintf("Invalid number of values for parameter %s", key))
		return "", false
	}
	if key == "id" {
		id, err := strconv.Atoi(values[0])
		if err != nil {
			log.Print(id)
			writeResponse(w, http.StatusBadRequest, err, fmt.Sprintf("%s is not an integer", values[0]))
			return "", false
		}
	}
	return values[0], true
}

func writeResponse(w http.ResponseWriter, statusCode int, err error, message string) {
	log.Printf("%s: %s", message, err)
	w.WriteHeader(statusCode)
}

func loadConfig(cfgFile string) *configuration {
	log.Printf("Loading configuration from: %s", cfgFile)

	bs, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		log.Fatalf("Failed to open file %s: %s", cfgFile, err)
	}

	config := new(configuration)
	if err := json.Unmarshal(bs, config); err != nil {
		log.Fatalf("Failed to unmarshal configuration file: %s", err)
	}

	if config.PORT <= 0 {
		log.Fatal("The port must be larger than 0")
	}

	return config
}
