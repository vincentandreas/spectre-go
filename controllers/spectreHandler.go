package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"spectre-go/cmd"
	"spectre-go/models"
)

type BaseHandler struct {
	repo models.SiteResultRepository
}

func HandleRequests(h *BaseHandler) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/api/generatePassword", h.ProcessGenPasswd).Methods("POST")

	router.HandleFunc("/api/health", h.CheckHealth).Methods("GET")

	router.HandleFunc("/", h.MainPage)
	return router
}

func NewBaseHandler(repo models.SiteResultRepository) *BaseHandler {
	return &BaseHandler{
		repo: repo,
	}
}

func (h *BaseHandler) CheckHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	temp := map[string]string{
		"result": "OK",
	}
	json.NewEncoder(w).Encode(temp)
}

func (h *BaseHandler) MainPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		//http.FileServer(http.Dir("static/"))
		t, _ := template.ParseFiles("static/index.html")
		t.Execute(w, nil)
	}
}

func (h *BaseHandler) ProcessGenPasswd(w http.ResponseWriter, r *http.Request) {
	log.Printf("--- Begin Process gen password ---")
	w.Header().Set("Content-Type", "application/json")
	utilizeCache := r.Header.Get("Utilize-Cache")

	reqBody, _ := ioutil.ReadAll(r.Body)
	var genParams models.GenSiteParam
	json.Unmarshal(reqBody, &genParams)

	validate := validator.New()
	err := validate.Struct(genParams)

	if err != nil {
		fmt.Println("Validation failed")
		w.WriteHeader(http.StatusBadRequest)
		temp := models.ApiResponse{"", "01", "Request not valid"}
		json.NewEncoder(w).Encode(temp)
		return
	}
	var hashedKey string
	if utilizeCache == "true" {
		hashedKey = cmd.HashParams(genParams)

		log.Printf("Finding from Cache...")
		cacheRes := h.repo.FindSiteResult(hashedKey)
		log.Printf("Raw Cache = " + cacheRes)
		cacheContent := cmd.DecryptContent(genParams.Username, genParams.Password, cacheRes)
		if cacheContent != "" {
			log.Printf("Cache found!")

			temp := models.ApiResponse{
				cacheContent,
				"00",
				"Success",
			}

			json.NewEncoder(w).Encode(temp)
			return
		}
	}

	genResult := cmd.NewSiteResult(genParams)

	if utilizeCache == "true" {
		log.Printf("Saving to cache")

		encResult := cmd.EncryptContent(genParams.Username, genParams.Password, genResult)
		log.Printf("Encrypted content : " + encResult)
		h.repo.Save(hashedKey, encResult)
	}
	temp := models.ApiResponse{
		genResult,
		"00",
		"Success",
	}

	json.NewEncoder(w).Encode(temp)
}
