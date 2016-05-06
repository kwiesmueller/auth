package ledis

import (
	"fmt"
	"net"
	"os"
	"testing"

	"path"

	. "github.com/bborbe/assert"
	lediscfg "github.com/siddontang/ledisdb/config"
	"github.com/siddontang/ledisdb/server"
)

func TestImplementsClient(t *testing.T) {
	object := New("", "")
	var expected *Client
	err := AssertThat(object, Implements(expected))
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetSet(t *testing.T) {
	port, err := getFreePort()
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	server, err := createServer(port)
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	go server.Run()

	client := New(fmt.Sprintf("localhost:%d", port), "secret")
	err = client.Set("hello", "world")
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	value, err := client.Get("hello")
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(value, Is("world")); err != nil {
		t.Fatal(err)
	}
}

func TestDel(t *testing.T) {
	port, err := getFreePort()
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	server, err := createServer(port)
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	go server.Run()

	client := New(fmt.Sprintf("localhost:%d", port), "secret")
	err = client.Set("hello", "world")
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	value, err := client.Get("hello")
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(value, Is("world")); err != nil {
		t.Fatal(err)
	}
	client.Del("hello")
	_, err = client.Get("hello")
	if err := AssertThat(err, NotNilValue()); err != nil {
		t.Fatal(err)
	}
}

func TestHashGetSet(t *testing.T) {
	port, err := getFreePort()
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	server, err := createServer(port)
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	go server.Run()

	client := New(fmt.Sprintf("localhost:%d", port), "secret")
	err = client.HashSet("hello", "new", "world")
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	value, err := client.HashGet("hello", "new")
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(value, Is("world")); err != nil {
		t.Fatal(err)
	}
}

func TestHashDel(t *testing.T) {
	port, err := getFreePort()
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	server, err := createServer(port)
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	go server.Run()

	client := New(fmt.Sprintf("localhost:%d", port), "secret")
	err = client.HashSet("hello", "new", "world")
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	value, err := client.HashGet("hello", "new")
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(value, Is("world")); err != nil {
		t.Fatal(err)
	}
	err = client.HashDel("hello", "new")
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	_, err = client.HashGet("hello", "new")
	if err := AssertThat(err, NotNilValue()); err != nil {
		t.Fatal(err)
	}
}

func TestHashClear(t *testing.T) {
	port, err := getFreePort()
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	server, err := createServer(port)
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	go server.Run()

	client := New(fmt.Sprintf("localhost:%d", port), "secret")
	err = client.HashSet("hello", "new", "world")
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	value, err := client.HashGet("hello", "new")
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(value, Is("world")); err != nil {
		t.Fatal(err)
	}
	err = client.HashClear("hello")
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	_, err = client.HashGet("hello", "new")
	if err := AssertThat(err, NotNilValue()); err != nil {
		t.Fatal(err)
	}
}

func TestPingSuccess(t *testing.T) {
	port, err := getFreePort()
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	server, err := createServer(port)
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	go server.Run()

	client := New(fmt.Sprintf("localhost:%d", port), "secret")
	err = client.Ping()
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
}

func TestPingFailure(t *testing.T) {
	client := New(fmt.Sprintf("localhost:1234"), "secret")
	err := client.Ping()
	if err := AssertThat(err, NotNilValue()); err != nil {
		t.Fatal(err)
	}
}

func createServer(port int) (*server.App, error) {
	cfg := lediscfg.NewConfigDefault()
	cfg.AuthPassword = "secret"
	cfg.Addr = fmt.Sprintf("localhost:%d", port)
	cfg.DBPath = path.Join(os.TempDir(), "/dbdir")
	cfg.DataDir = path.Join(os.TempDir(), "/datadir")
	return server.NewApp(cfg)
}

func getFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

func TestExists(t *testing.T) {
	port, err := getFreePort()
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	server, err := createServer(port)
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	go server.Run()

	client := New(fmt.Sprintf("localhost:%d", port), "secret")
	result, err := client.Exists("hello")
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(result, Is(false)); err != nil {
		t.Fatal(err)
	}
	err = client.Set("hello", "world")
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	result, err = client.Exists("hello")
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(result, Is(true)); err != nil {
		t.Fatal(err)
	}
}

func TestHashExists(t *testing.T) {
	port, err := getFreePort()
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	server, err := createServer(port)
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	go server.Run()

	client := New(fmt.Sprintf("localhost:%d", port), "secret")
	exists, err := client.HashExists("hello", "new")
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(exists, Is(false)); err != nil {
		t.Fatal(err)
	}
	err = client.HashSet("hello", "new", "world")
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	exists, err = client.HashExists("hello", "new")
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(exists, Is(true)); err != nil {
		t.Fatal(err)
	}
}
