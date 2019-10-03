package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/ronna-s/baby-janus/gateway"
)

const (
	domain = "http://127.0.0.1:8080"
	orig   = "/orig"
	dest   = "/dest"
)

func TestProxy(t *testing.T) {
	setup()

	testProxy(t)
}

func testProxy(t *testing.T) {
	rs := "good job"
	found := false

	cookie := http.Cookie{Name: "test", Value: "cookie"}
	http.HandleFunc(dest, func(w http.ResponseWriter, r *http.Request) {
		if len(r.Cookies()) == 0 || r.Cookies()[0].Name != cookie.Name || r.Cookies()[0].Value != cookie.Value {
			t.Fatalf("Expcetd to receive the original request's cookie")
		}
		found = true
		fmt.Fprintf(w, rs)
	})

	client := gateway.NewClient(domain)
	if err := client.RegisterRoute(orig, fmt.Sprintf("%s%s", domain, dest)); err != nil {
		t.Fatal(err)
	}

	body := string(httpGet(t, fmt.Sprintf("%s%s", domain, orig), &cookie))
	if rs != body {
		t.Error("http body is incorrect", body)
	}

	if !found {
		t.Error("failed to reach destination")
	}
}

func httpGet(t *testing.T, url string, cookie *http.Cookie) []byte {
	req, _ := http.NewRequest("GET", url, nil)
	if cookie != nil {
		req.AddCookie(cookie)
	}
	var client http.Client
	resp, err := client.Do(req)
	defer func() {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}()

	if resp.StatusCode == http.StatusNotFound {
		t.Fatalf("Not found %s", url)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected to return status ok. url: %s, returned: %d", url, resp.StatusCode)
	}

	if err != nil {
		t.Fatalf("unexpected error %s", err.Error())
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("unexpected error %s", err.Error())
	}
	return body
}

func setup() {
	ready := make(chan struct{})
	go run(ready)
	<-ready
}
