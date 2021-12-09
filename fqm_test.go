package main

import (
	"testing"
)

func TestNewFQm(t *testing.T) {
	fqm := NewFQm("./abc.qmcflac", "./")
	err := fqm.Decrypt()
	if err != nil {
		t.Fatalf("%s\n", err)
	}
}
