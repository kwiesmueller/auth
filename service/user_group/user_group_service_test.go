package user_group

import (
	"testing"

	. "github.com/bborbe/assert"
	"github.com/bborbe/auth/service"
)

func TestImplementsService(t *testing.T) {
	object := New(nil, nil)
	var expected *service.UserGroupService
	err := AssertThat(object, Implements(expected))
	if err != nil {
		t.Fatal(err)
	}
}
