package application

import (
	"testing"

	. "github.com/bborbe/assert"
	"github.com/bborbe/auth/service"
)

func TestImplementsService(t *testing.T) {
	object := New()
	var expected *service.ApplicationService
	if err := AssertThat(object, Implements(expected)); err != nil {
		t.Fatal(err)
	}
}
