package user

import (
	"testing"

	. "github.com/bborbe/assert"
	"github.com/bborbe/auth/service"
)

func TestImplementsService(t *testing.T) {
	object := New()
	var expected *service.UserService
	err := AssertThat(object, Implements(expected))
	if err != nil {
		t.Fatal(err)
	}
}
