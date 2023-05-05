package cmd

import (
	"fmt"
	"testing"
)

func TestHashParams(t *testing.T) {
	parameters := GenerateMedParams()

	for i := range parameters {
		t.Run(fmt.Sprintf("Testing [%v]", i), func(t *testing.T) {
			actual := HashParams(parameters[i].input)
			if actual != parameters[i].expectedHash {
				t.Logf("expectedPasswd : %s: , actual: %s", parameters[i].expectedPasswd, actual)
				t.Fail()
			}
		})
	}
}
