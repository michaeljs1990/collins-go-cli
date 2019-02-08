// Helper functions I pulled out of the go-collins library
// to help with testing this client.
package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	collins "gopkg.in/tumblr/go-collins.v0/collins"
)

var (
	mux    *http.ServeMux
	client *collins.Client
	server *httptest.Server
)

func setup() *collins.Client {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	client, _ = collins.NewClient("test", "test", "test")
	url, _ := url.Parse(server.URL)
	client.BaseURL = url

	return client
}

func teardown() {
	server.Close()
}

func SetupMethod(code int, method, url, file string, t *testing.T) {
	mux.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			t.Errorf("Request method: %v, want %s", r.Method, method)
		}
		resp, err := ioutil.ReadFile(file)
		if err != nil {
			t.Errorf("Could not read %s\n", file)
		}
		w.Header().Set("Content-Type", "application/json;")
		w.WriteHeader(code)
		fmt.Fprintf(w, "%s", resp)
	})
}

func SetupGET(code int, url, file string, t *testing.T) {
	SetupMethod(code, "GET", url, file, t)
}

func SetupPUT(code int, url, file string, t *testing.T) {
	SetupMethod(code, "PUT", url, file, t)
}

func SetupPOST(code int, url, file string, t *testing.T) {
	SetupMethod(code, "POST", url, file, t)
}

func SetupDELETE(code int, url, file string, t *testing.T) {
	SetupMethod(code, "DELETE", url, file, t)
}

func captureStdout() (chan string, *os.File, *os.File) {
	old := os.Stdout
	r, w, _ := os.Pipe()

	output := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		output <- buf.String()
	}()

	os.Stdout = w

	return output, old, w
}

func returnStdout(c chan string, old *os.File, w *os.File) string {
	w.Close()
	os.Stdout = old
	return strings.Trim(<-c, "\n")
}
