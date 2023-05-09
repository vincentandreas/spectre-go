package models

import (
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenSiteParam_failed_when_username_blank(t *testing.T) {
	req := GenSiteParam{
		Username:   "",
		Password:   "qww",
		KeyPurpose: "password",
		Site:       "twitter.com",
		KeyType:    "phrase",
		KeyCounter: 1,
	}

	val := validator.New()
	err := val.Struct(req)
	validationErr := err.(validator.ValidationErrors)
	assert.Contains(t, validationErr[0].Error(), "'Username' failed")
}

func TestGenSiteParam_failed_when_password_blank(t *testing.T) {
	req := GenSiteParam{
		Username:   "qwww",
		Password:   "",
		KeyPurpose: "password",
		Site:       "twitter.com",
		KeyType:    "phrase",
		KeyCounter: 1,
	}

	val := validator.New()
	err := val.Struct(req)
	validationErr := err.(validator.ValidationErrors)
	assert.Contains(t, validationErr[0].Error(), "'Password' failed")
}

func TestGenSiteParam_failed_when_purpose_unknown(t *testing.T) {
	req := GenSiteParam{
		Username:   "qwww",
		Password:   "qwww",
		KeyPurpose: "testingfailed",
		Site:       "twitter.com",
		KeyType:    "phrase",
		KeyCounter: 1,
	}

	val := validator.New()
	err := val.Struct(req)
	validationErr := err.(validator.ValidationErrors)
	assert.Contains(t, validationErr[0].Error(), "'KeyPurpose' failed")
}

func TestGenSiteParam_failed_when_site_empty(t *testing.T) {
	req := GenSiteParam{
		Username:   "qwww",
		Password:   "qwww",
		KeyPurpose: "password",
		Site:       "",
		KeyType:    "phrase",
		KeyCounter: 1,
	}

	val := validator.New()
	err := val.Struct(req)
	validationErr := err.(validator.ValidationErrors)
	assert.Contains(t, validationErr[0].Error(), "'Site' failed")
}

func TestGenSiteParam_failed_when_counter_invalid_num(t *testing.T) {
	req := GenSiteParam{
		Username:   "qwww",
		Password:   "qwww",
		KeyPurpose: "password",
		Site:       "twitter.com",
		KeyType:    "med",
		KeyCounter: -10,
	}

	val := validator.New()
	err := val.Struct(req)
	validationErr := err.(validator.ValidationErrors)
	assert.Contains(t, validationErr[0].Error(), "'KeyCounter' failed")
}

func TestGenSiteParam_success_when_request_valid(t *testing.T) {
	req := GenSiteParam{
		Username:   "qwww",
		Password:   "qwww",
		KeyPurpose: "password",
		Site:       "twitter.com",
		KeyType:    "med",
		KeyCounter: 1,
	}

	val := validator.New()
	err := val.Struct(req)
	assert.Nil(t, err)
}
