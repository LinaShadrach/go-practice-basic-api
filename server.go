package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"os"

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

	if id := r.URL.Query().Get("id"); id != "" {
		handleIDRequest(w, r, id)
	} else {
		handleCountRequest(w, r)
	}
}

func handleCountRequest(w http.ResponseWriter, r *http.Request) {
	all, err := dbconn.Query("select count from viewers")
	if err != nil {
		log.Printf("Failed to query the DB: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer all.Close()

	var count, total int
	for all.Next() {
		if err := all.Scan(&count); err != nil {
			log.Printf("Error while iterating over rows: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		total += count
	}
	if hadError(w, all) {
		return
	}

	fmt.Fprintf(w, "%d", total)
}

func handleIDRequest(w http.ResponseWriter, r *http.Request, idstr string) {
	// validate the users input
	id, err := strconv.Atoi(idstr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Invalid id parameter, must be an integer")
		return
	}

	// request the count for that entry
	rows, err := dbconn.Query("select count from viewers where id = $1", id)
	if err != nil {
		log.Printf("Failed to query the DB: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// check if we found anything
	if rows.Next() {
		var count int
		if err := rows.Scan(&count); err != nil {
			log.Printf("Error while iterating over rows: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if hadError(w, rows) {
			return
		}

		fmt.Fprintf(w, "%d", count)
		return
	}
	if hadError(w, rows) {
		return
	}
	w.WriteHeader(http.StatusNotFound)
}

func hadError(w http.ResponseWriter, rows *sql.Rows) bool {
	if err := rows.Err(); err != nil {
		log.Printf("Error while interacting with the DB: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return true
	}
	return false
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
