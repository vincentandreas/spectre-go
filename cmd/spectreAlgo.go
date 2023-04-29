package cmd

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"golang.org/x/crypto/scrypt"
	"log"
)

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
	fmt.Println("Key result : ")
	for _, data := range userKeyData {
		fmt.Print(data, " ")
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
	fmt.Println("----------------------------")
	for _, data := range keyData {
		fmt.Print(data, ".")
	}
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

func newSiteResult(username string, password string, site string, keyCounter int, keyPurpose string, keyType string) string {
	userKey := newUserKey(username, password, keyPurpose)
	siteKey := newSiteKey(userKey, site, keyCounter, keyPurpose, "")
	resTemplates := templates[keyType]
	resTemplate := resTemplates[int(siteKey[0])%len(resTemplates)]
	passRes := ""
	for i := 0; i < len(resTemplate); i++ {
		currChar := characters[string(resTemplate[i])]
		idx := int(siteKey[i+1]) % len(currChar)
		passRes += string([]rune(currChar)[idx])
	}
	fmt.Println("pass ress = " + passRes)
	return passRes
}
