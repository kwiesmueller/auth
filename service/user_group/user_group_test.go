package user_group

import (
	"testing"

	. "github.com/bborbe/assert"
)

func TestImplementsService(t *testing.T) {
	object := New(nil, nil)
	var expected *Service
	err := AssertThat(object, Implements(expected))
	if err != nil {
		t.Fatal(err)
	}
}
