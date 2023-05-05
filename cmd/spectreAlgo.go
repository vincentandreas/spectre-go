package cmd

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
	"golang.org/x/crypto/scrypt"
	"log"
	"spectre-go/models"
	"strings"
)

var tempCache = make(map[string]string)

func convBigEndian(numb int) []byte {
	numbByte := make([]byte, 4)
	binary.BigEndian.PutUint32(numbByte, uint32(numb))
	return numbByte
}

func newUserKey(username string, password string, purpose string) []byte {
	usernameBytes := []byte(username)
	passwordBytes := []byte(password)
	purposeBytes := []byte(purpose)
	saltLgth := len(purposeBytes) + 4 + len(usernameBytes)
	//var userSalt []byte
	userSalt := make([]byte, saltLgth)
	uS := 0
	for ctr, data := range purposeBytes {
		userSalt[ctr] = data
	}
	uS += len(purposeBytes)

	for ctr, data := range convBigEndian(len(usernameBytes)) {
		userSalt[ctr+uS] = data
	}
	uS += 4
	for ctr, data := range usernameBytes {
		userSalt[ctr+uS] = data
	}
	userKeyData, err := scrypt.Key(passwordBytes, userSalt, 32768, 8, 2, 64)
	if err != nil {
		log.Fatal(err)
	}
	return userKeyData
}

func newSiteKey(userKeyCrypto []byte, siteName string, keyCounter int, purpose string, keyContext string) []byte {
	siteNameBytes := []byte(siteName)
	purposeBytes := []byte(purpose)

	keyContextLen := 0
	if keyContext != "" {
		keyContextBytes := []byte(keyContext)
		keyContextLen = len(keyContextBytes)
	}
	saltLgth := len(purposeBytes) + 4 + len(siteNameBytes) + 4 + keyContextLen
	siteSalt := make([]byte, saltLgth)
	sS := 0
	for ctr, data := range purposeBytes {
		siteSalt[ctr] = data
	}
	sS += len(purposeBytes)

	for ctr, data := range convBigEndian(len(siteNameBytes)) {
		siteSalt[ctr+sS] = data
	}
	sS += 4

	for ctr, data := range siteNameBytes {
		siteSalt[ctr+sS] = data
	}
	sS += len(siteNameBytes)

	for ctr, data := range convBigEndian(keyCounter) {
		siteSalt[ctr+sS] = data
	}
	sS += 4

	h := hmac.New(sha256.New, userKeyCrypto)

	// Write Data to it
	h.Write(siteSalt)

	// Get result and encode as hexadecimal string
	keyData := h.Sum(nil)
	return keyData
}

var templates = map[string][]string{
	"med": {"CvcnoCvc", "CvcCvcno"},
	"long": {
		"CvcvnoCvcvCvcv",
		"CvcvCvcvnoCvcv",
		"CvcvCvcvCvcvno",
		"CvccnoCvcvCvcv",
		"CvccCvcvnoCvcv",
		"CvccCvcvCvcvno",
		"CvcvnoCvccCvcv",
		"CvcvCvccnoCvcv",
		"CvcvCvccCvcvno",
		"CvcvnoCvcvCvcc",
		"CvcvCvcvnoCvcc",
		"CvcvCvcvCvccno",
		"CvccnoCvccCvcv",
		"CvccCvccnoCvcv",
		"CvccCvccCvcvno",
		"CvcvnoCvccCvcc",
		"CvcvCvccnoCvcc",
		"CvcvCvccCvccno",
		"CvccnoCvcvCvcc",
		"CvccCvcvnoCvcc",
		"CvccCvcvCvccno",
	},
	"max": {
		"anoxxxxxxxxxxxxxxxxx",
		"axxxxxxxxxxxxxxxxxno",
	},
	"short": {
		"Cvcn",
	},
	"basic": {
		"aaanaaan",
		"aannaaan",
		"aaannaaa",
	},
	"pin": {
		"nnnn",
	},
	"name": {
		"cvccvcvcv",
	},
	"phrase": {
		"cvcc cvc cvccvcv cvc",
		"cvc cvccvcvcv cvcv",
		"cv cvccv cvc cvcvccv",
	},
}
var characters = map[string]string{
	"V": "AEIOU",
	"C": "BCDFGHJKLMNPQRSTVWXYZ",
	"v": "aeiou",
	"c": "bcdfghjklmnpqrstvwxyz",
	"A": "AEIOUBCDFGHJKLMNPQRSTVWXYZ",
	"a": "AEIOUaeiouBCDFGHJKLMNPQRSTVWXYZbcdfghjklmnpqrstvwxyz",
	"n": "0123456789",
	"o": "@&%?,=[]_:-+*$#!'^~;()/.",
	"x": "AEIOUaeiouBCDFGHJKLMNPQRSTVWXYZbcdfghjklmnpqrstvwxyz0123456789!@#$%^&*()",
	" ": " ",
}

func NewSiteResult(params models.GenSiteParam) string {
	log.Printf("Running NewSiteResult")

	userKey := newUserKey(params.Username, params.Password, params.KeyPurpose)
	siteKey := newSiteKey(userKey, params.Site, params.KeyCounter, params.KeyPurpose, "")
	resTemplates := templates[params.KeyType]
	resTemplate := resTemplates[int(siteKey[0])%len(resTemplates)]
	var passRes strings.Builder
	for i := 0; i < len(resTemplate); i++ {
		currChar := characters[string(resTemplate[i])]
		idx := int(siteKey[i+1]) % len(currChar)
		passRes.WriteRune([]rune(currChar)[idx])
	}

	return passRes.String()
}
