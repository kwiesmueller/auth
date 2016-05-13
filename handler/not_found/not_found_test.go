package not_found

import (
	"net/http"
	"testing"

	. "github.com/bborbe/assert"
)

func TestImplementsHandler(t *testing.T) {
	object := New()
	var expected *http.Handler
	err := AssertThat(object, Implements(expected))
	if err != nil {
		t.Fatal(err)
	}
}
