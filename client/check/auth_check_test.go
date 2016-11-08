package check

import (
	"testing"

	. "github.com/bborbe/assert"
)

func TestNew(t *testing.T) {
	c := New(nil)
	err := AssertThat(c, NotNilValue())
	if err != nil {
		t.Fatal(err)
	}
}
