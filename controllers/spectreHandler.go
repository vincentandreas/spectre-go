package controllers

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
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

	router.HandleFunc("/api/getPassword", h.ProcessGenPasswd).Methods("POST")

	router.HandleFunc("/api/health", h.CheckHealth).Methods("GET")

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
		panic("Validation failed")
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
			temp := map[string]string{
				"result": cacheContent,
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

	temp := map[string]string{
		"result": genResult,
	}
	json.NewEncoder(w).Encode(temp)
}
