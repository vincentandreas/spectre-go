package cmd

import (
	"crypto/sha256"
	"encoding/base64"
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
