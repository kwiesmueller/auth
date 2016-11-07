package main

import (
	"testing"

	. "github.com/bborbe/assert"
)

func TestServerConfig(t *testing.T) {
	config := createConfig()
	if err := AssertThat(config, NotNilValue()); err != nil {
		t.Fatal(err)
	}
}
