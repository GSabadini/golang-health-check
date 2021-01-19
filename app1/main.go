package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	//Gorilla router
	router := mux.NewRouter()

	router.HandleFunc("/health", Health)

	log.Println("Start HTTP server :3001")
	if err := http.ListenAndServe(":3001", router); err != nil {
		panic(err)
	}
}

func Health(w http.ResponseWriter, _ *http.Request) {
	var response = map[string]interface{}{
		"message":     "OK",
		"http_status": http.StatusOK,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
