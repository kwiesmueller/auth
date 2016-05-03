package user_finder

import (
	"testing"

	. "github.com/bborbe/assert"
)

func TestImplementsUserFinder(t *testing.T) {
	object := New()
	var expected *UserFinder
	err := AssertThat(object, Implements(expected))
	if err != nil {
		t.Fatal(err)
	}
}
