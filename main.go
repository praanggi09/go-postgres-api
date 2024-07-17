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
	VmName           string `json:"vmName"`
	Description      string `json:"description"`
	OperatingSystem  string `json:"os"`
	Username         string `json:"username"`
	Password         string `json:"password"`
	IP               string `json:"ip"`
	Hostname         string `json:"hostname"`
	ProvisionedSpace int    `json:"provisionedSpace"`
	UsedSpace        int    `json:"usedSpace"`
	MemorySize       int    `json:"memorySize"`
	CPU              int    `json:"cpu"`
	Notes            string `json:"notes"`
}

func initDB() {
	var err error
	connStr := "user=postgres dbname=effortease password=postgresP@ssw0rd!234 host=10.87.1.122 port=5432 sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println("Connection Succes")
	defer db.Close()

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

func getResponses(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT * FROM SERVERS")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	responses := []Response{}
	for rows.Next() {
		var response Response
		if err := rows.Scan(&response.ID, &response.ResourcePool); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		responses = append(responses, response)
	}

	json.NewEncoder(w).Encode(responses)
}

func createResponse(w http.ResponseWriter, r *http.Request) {
	var response Response
	if err := json.NewDecoder(r.Body).Decode(&response); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := db.QueryRow("INSERT INTO SERVERS (resourcePool) VALUES ($1) RETURNING ID", response.ResourcePool).Scan(&response.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(response)
}
