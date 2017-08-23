package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/wwgberlin/baby_janus/gateway/gateway-api"
)

const (
	domain = "http://127.0.0.1:8080"
	orig   = "/orig"
	dest   = "/dest"
)

func TestRedirects(t *testing.T) {
	setup()

	t.Run("Test Redirect - Simple", testRedirect)
}

func testRedirect(t *testing.T) {

	rs := "good job"
	rc := false

	http.HandleFunc(dest, func(w http.ResponseWriter, r *http.Request) {
		rc = true
		fmt.Fprintf(w, rs)
	})

	client := gateway_api.NewGatewayClient(domain)
	client.RegisterRoute(orig, fmt.Sprintf("%s%s", domain, dest))

	body := string(mustHttpGet(t, fmt.Sprintf("%s%s", domain, orig)))
	if rs != body {
		t.Error("http body is incorrect", body)
	}

	if !rc {
		t.Error("failed to reach destination")
	}
}

func mustHttpGet(t *testing.T, url string) (body []byte) {
	resp, err := http.Get(url)
	defer func() {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}()

	if resp.StatusCode == http.StatusNotFound {
		panic(fmt.Sprintf("Not found %s", url))
	}

	if resp.StatusCode != http.StatusOK {
		panic(fmt.Sprintf("Expected to return status ok. url: %s, returned: %s", url, resp.StatusCode))
	}

	if err != nil {
		panic(err)
	}
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return body
}

func setup() {
	go main()
	<-time.After(10 * time.Millisecond) //give the server some time to start
}
