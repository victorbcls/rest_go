package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"rest/db"
	"rest/response"
	"rest/users"

	"github.com/gorilla/mux"
)

func getUsers(w http.ResponseWriter, r *http.Request) {
	users, err := users.GetUsers(Client)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response.NewResponseHttp(500, "Server side error"))
	} else {
		json.NewEncoder(w).Encode(users)

	}
}

func getUserByCpf(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var user users.User
	user, err := users.GetUserByCpf(Client, params["cpf"])
	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response.NewResponseHttp(404, "User not found"))
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response.NewResponseHttp(500, "Server side error"))
	} else {
		json.NewEncoder(w).Encode(response.NewDataResponse(200, user))

	}

}
func createUser(w http.ResponseWriter, r *http.Request) {
	var user users.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response.NewResponseHttp(500, "Server side error"))

	}
	response := users.CreateUser(Client, user)
	w.WriteHeader(response.Status)
	json.NewEncoder(w).Encode(response)

}

var Client *sql.DB

func main() {

	Client, _ = db.Connect()
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/users", getUsers).Methods("GET")
	router.HandleFunc("/users/{cpf}", getUserByCpf).Methods("GET")
	router.HandleFunc("/users/create", createUser).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", router))
}
