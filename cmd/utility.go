package cmd

import (
	"bufio"
	"os"
	"spectre-go/models"
	"strconv"
	"strings"
)

type SingleIe struct {
	input          models.GenSiteParam
	expectedPasswd string
	expectedHash   string
}

func GenerateMedParams() []SingleIe {
	readFile, err := os.Open("testcases.txt")
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(readFile)
	scanner.Split(bufio.ScanLines)
	var parameters []SingleIe

	for scanner.Scan() {
		singleLine := scanner.Text()
		inputSplit := strings.Split(singleLine, ",")

		parsedKeyCtr, _ := strconv.ParseInt(inputSplit[4], 10, 64)

		temp := SingleIe{
			models.GenSiteParam{
				Username:   inputSplit[0],
				Password:   inputSplit[1],
				Site:       inputSplit[2],
				KeyPurpose: inputSplit[3],
				KeyCounter: int(parsedKeyCtr),
				KeyType:    inputSplit[5],
			},
			inputSplit[6],
			inputSplit[7],
		}
		parameters = append(parameters, temp)

	}
	return parameters
}
