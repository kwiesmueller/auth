package version

import (
	"net/http"
	"testing"

	"os"

	. "github.com/bborbe/assert"
	"github.com/golang/glog"
)

func TestMain(m *testing.M) {
	exit := m.Run()
	glog.Flush()
	os.Exit(exit)
}

func TestImplementsHandler(t *testing.T) {
	object := New()
	var expected *http.Handler
	if err := AssertThat(object, Implements(expected)); err != nil {
		t.Fatal(err)
	}
}
