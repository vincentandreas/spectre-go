package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"testing"
)

func BenchmarkNewSiteResult(b *testing.B) {
	parameters := GenerateMedParams()
	log.SetOutput(ioutil.Discard)

	for x := range parameters {
		b.Run(fmt.Sprintf("Benchmark [%v]", x), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				NewSiteResult(parameters[x].input)
			}
		})
	}
}

func TestNewSiteResultParameters(t *testing.T) {
	parameters := GenerateMedParams()

	for i := range parameters {
		t.Run(fmt.Sprintf("Testing [%v]", i), func(t *testing.T) {
			actual := NewSiteResult(parameters[i].input)
			if actual != parameters[i].expectedPasswd {
				t.Logf("expectedPasswd : %s: , actual: %s", parameters[i].expectedPasswd, actual)
				t.Fail()
			}
		})
	}
}
