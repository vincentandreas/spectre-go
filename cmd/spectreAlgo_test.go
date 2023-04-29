package cmd

import "testing"

func TestNewSiteResult_med(t *testing.T) {
	params := GenSiteParam{
		Username:   "a",
		Password:   "a",
		Site:       "twitter.com",
		KeyPurpose: "com.lyndir.masterpassword",
		KeyCounter: 1,
		KeyType:    "med",
	}
	res := NewSiteResult(params)
	if res != "RevXep5+" {
		t.Error("Password not same")
	}
}

func TestNewSiteResult_long(t *testing.T) {
	params := GenSiteParam{
		Username:   "a",
		Password:   "",
		Site:       "twitter.com",
		KeyPurpose: "com.lyndir.masterpassword",
		KeyCounter: 1,
		KeyType:    "long",
	}
	res := NewSiteResult(params)
	if res != "RevoGupsWunl3-" {
		t.Error("Password not same")
	}
}
