package cmd

import (
	"bufio"
	"fmt"
	"os"
	"spectre-go/models"
	"strconv"
	"strings"
	"testing"
)

func BenchmarkNewSiteResult(b *testing.B) {
	parameters := generateMedParams()

	for x := range parameters {
		b.Run(fmt.Sprintf("Benchmark [%v]", x), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				NewSiteResult(parameters[x].input)
			}
		})
	}
}

func TestNewSiteResultParameters(t *testing.T) {
	parameters := generateMedParams()

	for i := range parameters {
		t.Run(fmt.Sprintf("Testing [%v]", i), func(t *testing.T) {
			actual := NewSiteResult(parameters[i].input)
			if actual != parameters[i].expected {
				t.Logf("expected : %s: , actual: %s", parameters[i].expected, actual)
				t.Fail()
			}
		})
	}
}

type SingleIe struct {
	input    models.GenSiteParam
	expected string
}

func generateMedParams() []SingleIe {
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
			}, inputSplit[6],
		}
		parameters = append(parameters, temp)

	}
	return parameters
}
