package handler_creator

import (
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

func TestImplementsHandlerCreator(t *testing.T) {
	object := New()
	var expected *HandlerCreator
	if err := AssertThat(object, Implements(expected)); err != nil {
		t.Fatal(err)
	}
}
