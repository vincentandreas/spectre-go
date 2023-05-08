package cmd

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"golang.org/x/crypto/scrypt"
	"io"
	"log"
	"spectre-go/models"
	"strconv"
)

func HashParams(param models.GenSiteParam) string {
	h := sha256.New()
	keys := "uname:" + param.Username + "|passwd:" + param.Password +
		"|site:" + param.Site + "|keyCtr:" + strconv.Itoa(param.KeyCounter) + "|keyPurpose:" + param.KeyPurpose + "|keyType:" + param.KeyType
	h.Write([]byte(keys))
	sha := base64.URLEncoding.EncodeToString(h.Sum(nil))
	return sha
}

func EncryptContent(username string, password string, content string) string {
	scryptKey, err := generateScryptKey(username, password)
	if err != nil {
		log.Fatal(err)
	}
	c, err := aes.NewCipher(scryptKey)
	if err != nil {
		log.Fatal(err)
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		log.Fatal(err)
	}
	nonce := make([]byte, gcm.NonceSize())

	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		log.Fatal(err)
	}

	encryptRes := gcm.Seal(nonce, nonce, []byte(content), nil)
	encoded := base64.URLEncoding.EncodeToString(encryptRes)
	return encoded
}

func DecryptContent(username string, password string, encryptedContent string) string {
	if len(encryptedContent) == 0 {
		return ""
	}
	ciphertext, err := base64.URLEncoding.DecodeString(encryptedContent)
	if err != nil {
		log.Print(err)
	}

	scryptKey, err := generateScryptKey(username, password)
	c, err := aes.NewCipher(scryptKey)
	if err != nil {
		log.Print(err)
	}

	gcm, err := cipher.NewGCM(c)

	if err != nil {
		log.Print(err)
	}

	nonceSize := gcm.NonceSize()
	if len(encryptedContent) < nonceSize {
		log.Print(err)
	}
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plainText, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		log.Print(err)
	}
	return string(plainText)
}

func generateScryptKey(username string, password string) ([]byte, error) {
	keys := "uname:" + username + "|passwd:" + password
	salt := "hulala"
	scryptKey, err := scrypt.Key([]byte(keys), []byte(salt), 32768, 8, 2, 32)
	return scryptKey, err
}
