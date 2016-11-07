package user_data

import (
	"testing"

	. "github.com/bborbe/assert"
	"github.com/bborbe/auth/service"
)

func TestImplementsService(t *testing.T) {
	object := New(nil)
	var expected *service.UserDataService
	err := AssertThat(object, Implements(expected))
	if err != nil {
		t.Fatal(err)
	}
}
