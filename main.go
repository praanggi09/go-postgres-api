package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var db *sql.DB

type Response struct {
	ID               string `json:"id"`
	ResourcePool     string `json:"resourcepool"`
	VmName           string `json:"vmname"`
	Description      string `json:"description"`
	OperatingSystem  string `json:"os"`
	Username         string `json:"username"`
	Password         string `json:"password"`
	IP               string `json:"ip"`
	Hostname         string `json:"hostname"`
	ProvisionedSpace int    `json:"provisionedspace"`
	UsedSpace        int    `json:"usedSpace"`
	MemorySize       int    `json:"memorySize"`
	CPU              int    `json:"cpu"`
	Notes            string `json:"notes"`
}

func initDB() {
	var err error
	// connStr := "user=postgres dbname=effortease password=oioi host=10.87.1.122 port=5432 sslmode=disable"
	connStr := "user=postgres dbname=golangDB password=P@ssw0rd!234 host=localhost port=5432 sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	initDB()

	router := mux.NewRouter()

	router.HandleFunc("/responses", getResponses).Methods("GET")
	router.HandleFunc("/responses", createResponse).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func getResponses(ww http.ResponseWriter, rr *http.Request) {
	rows, err := db.Query("SELECT id, resourcepool, vmname, description, os, username, password, ip, hostname, provisionedSpace, usedSpace, memorySize, cpu, notes FROM SERVERS")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var responses []Response
	for rows.Next() {
		var response Response
		if err := rows.Scan(&response.ID, &response.ResourcePool, &response.VmName, &response.Description, &response.OperatingSystem, &response.Username, &response.Password, &response.IP, &response.Hostname, &response.ProvisionedSpace, &response.UsedSpace, &response.MemorySize, &response.CPU, &response.Notes); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		responses = append(responses, response)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responses)
}

func createResponse(w http.ResponseWriter, r *http.Request) {
	var response Response
	if err := json.NewDecoder(r.Body).Decode(&response); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := db.QueryRow(
		"INSERT INTO SERVERS (resourcepool, vmName, description, os, username, password, ip, hostname, provisionedSpace, usedSpace, memorySize, cpu, notes) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) RETURNING id",
		response.ResourcePool, response.VmName, response.Description, response.OperatingSystem, response.Username, response.Password, response.IP, response.Hostname, response.ProvisionedSpace, response.UsedSpace, response.MemorySize, response.CPU, response.Notes,
	).Scan(&response.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func getResponse(w http.ResponseWriter, r * )