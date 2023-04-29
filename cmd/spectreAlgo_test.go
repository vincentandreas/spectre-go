package cmd

import "testing"

func TestNewSiteResult_med(t *testing.T) {
	res := newSiteResult("a", "a", "twitter.com", 1, "com.lyndir.masterpassword", "med")
	if res != "RevXep5+" {
		t.Error("Password not same")
	}
}

func TestNewSiteResult_long(t *testing.T) {
	res := newSiteResult("a", "a", "twitter.com", 1, "com.lyndir.masterpassword", "long")
	if res != "RevoGupsWunl3-" {
		t.Error("Password not same")
	}
}
