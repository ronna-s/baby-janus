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
	const (
		rs      = "good job"
		sup     = "sup?"
		method  = "YOLO"
		notMuch = "not much"
		retCode = http.StatusAlreadyReported
	)
	found := false

	cookie := http.Cookie{Name: "test", Value: "cookie"}
	http.HandleFunc(dest, func(w http.ResponseWriter, r *http.Request) {
		if len(r.Cookies()) == 0 || r.Cookies()[0].Name != cookie.Name || r.Cookies()[0].Value != cookie.Value {
			t.Errorf("Expcetd to receive the original request's cookie")
		}
		if r.Method != method {
			t.Errorf("Unexpected method in request")
		}
		if r.Header.Get("Content-Type") != sup {
			t.Errorf("Unexpected content type in request")
		}
		found = true
		w.Header().Add("Content-Type", notMuch)
		w.WriteHeader(retCode)
		fmt.Fprintf(w, rs)
	})

	client := gateway.NewClient(domain)
	if err := client.RegisterRoute(orig, fmt.Sprintf("%s%s", domain, dest)); err != nil {
		t.Fatal(err)
	}

	body, code, contentType := httpDo(method, t, fmt.Sprintf("%s%s", domain, orig), &cookie, sup)
	if rs != string(body) {
		t.Error("http body is incorrect", body)
	}
	if code != retCode {
		t.Errorf("Unexpcted response status code")
	}

	if contentType != notMuch {
		t.Errorf("Unexpected content type in response")
	}

	if !found {
		t.Error("failed to reach destination")
	}
}

func httpDo(method string, t *testing.T, url string, cookie *http.Cookie, contentType string) ([]byte, int, string) {
	req, _ := http.NewRequest(method, url, nil)
	if cookie != nil {
		req.AddCookie(cookie)
	}
	req.Header.Add("Content-Type", contentType)
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

	if err != nil {
		t.Fatalf("unexpected error %s", err.Error())
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("unexpected error %s", err.Error())
	}
	return body, resp.StatusCode, resp.Header.Get("Content-Type")
}

func setup() {
	ready := make(chan struct{})
	go run(ready)
	<-ready
}

func TestLoadBalancing(t *testing.T) {
	http.HandleFunc(dest+"1", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello from 1"))
	})
	http.HandleFunc(dest+"2", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello from 2"))
	})

	client := gateway.NewClient(domain)
	if err := client.RegisterRoute(orig+"new", fmt.Sprintf("%s%s", domain, dest+"1")); err != nil {
		t.Fatal(err)
	}
	if err := client.RegisterRoute(orig+"new", fmt.Sprintf("%s%s", domain, dest+"2")); err != nil {
		t.Fatal(err)
	}

	body, _, _ := httpDo("GET", t, fmt.Sprintf("%s%s", domain, orig+"new"), nil, "text/html")
	if "hello from 1" != string(body) {
		t.Fatal("http body is incorrect:", string(body))
	}

	body, _, _ = httpDo("GET", t, fmt.Sprintf("%s%s", domain, orig+"new"), nil, "text/html")
	if "hello from 2" != string(body) {
		t.Fatal("http body is incorrect:", string(body))
	}
}
