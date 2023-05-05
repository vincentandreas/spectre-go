package controllers

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"io/ioutil"
	"log"
	"net/http"
	"spectre-go/cmd"
	"spectre-go/models"
)

type BaseHandler struct {
	repo models.SiteResultRepository
}

func NewBaseHandler(repo models.SiteResultRepository) *BaseHandler {
	return &BaseHandler{
		repo: repo,
	}
}

func (h *BaseHandler) ProcessGenPasswd(w http.ResponseWriter, r *http.Request) {
	log.Printf("--- Begin Process gen password ---")
	w.Header().Set("Content-Type", "application/json")

	reqBody, _ := ioutil.ReadAll(r.Body)
	var genParams models.GenSiteParam
	json.Unmarshal(reqBody, &genParams)

	validate := validator.New()
	err := validate.Struct(genParams)

	if err != nil {
		panic("Validation failed")
	}

	hashedKey := cmd.HashParams(genParams)

	log.Printf("Finding from Cache...")
	cacheRes := h.repo.FindSiteResult(hashedKey)

	if cacheRes != "" {
		log.Printf("Cache found!")
		temp := map[string]string{
			"result": cacheRes,
		}
		json.NewEncoder(w).Encode(temp)
		return
	}

	genResult := cmd.NewSiteResult(genParams)

	log.Printf("Saving to cache")
	h.repo.Save(hashedKey, genResult)

	temp := map[string]string{
		"result": genResult,
	}
	json.NewEncoder(w).Encode(temp)
}
