package cmd

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHashParams(t *testing.T) {
	parameters := GenerateMedParams()

	for i := range parameters {
		t.Run(fmt.Sprintf("Testing [%v]", i), func(t *testing.T) {
			actual := HashParams(parameters[i].input)
			if actual != parameters[i].expectedHash {
				t.Logf("expectedHash : %s , actual: %s", parameters[i].expectedHash, actual)
				t.Fail()
			}
		})
	}
}

func TestDecrypt_should_empty_when_key_not_match(t *testing.T) {
	uname := "xxxxx"
	passwd := "xxxxx"
	encContent := "1_Vurqzh_gCsSHxph7uMFVhntMlPJBf5xN_0OxTSQeRu-RYh"
	decryptRes := DecryptContent(uname, passwd, encContent)
	fmt.Println(decryptRes)
	assert.Equal(t, decryptRes, "")
}

func TestEncryptDecrypt(t *testing.T) {
	uname := "grji0"
	passwd := "7xnbi1"
	content := "51n4bun6"

	encryptResult := EncryptContent(uname, passwd, content)
	fmt.Println("Encrypt result : " + encryptResult)
	decryptRes := DecryptContent(uname, passwd, encryptResult)
	fmt.Println("Decrypt result : " + decryptRes)

	assert.Equal(t, content, decryptRes)

}
