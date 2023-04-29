package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"spectre-go/cmd"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func processGenPasswd(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var genParams cmd.GenSiteParam
	json.Unmarshal(reqBody, &genParams)

	genResult := cmd.NewSiteResult(genParams)
	temp := map[string]string{
		"result": genResult,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(temp)
}

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/api/getPassword", processGenPasswd).Methods("POST")
	log.Fatal(http.ListenAndServe(":10000", router))
}

func main() {
	handleRequests()
}
