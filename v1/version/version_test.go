package version

import (
	"net/http"
	"testing"

	"os"

	"bytes"
	"fmt"

	. "github.com/bborbe/assert"
	"github.com/bborbe/io/writer_nop_close"
	"github.com/bborbe/log"
)

func TestMain(m *testing.M) {
	buffer := bytes.NewBufferString("")
	log.DefaultLogger = log.NewLogger(writer_nop_close.New(buffer))
	exit := m.Run()
	if exit != 0 {
		logger.Close()
		fmt.Print(buffer.String())
	}
	os.Exit(exit)
}

func TestImplementsHandler(t *testing.T) {
	object := New()
	var expected *http.Handler
	if err := AssertThat(object, Implements(expected)); err != nil {
		t.Fatal(err)
	}
}
