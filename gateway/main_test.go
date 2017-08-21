package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/wwgberlin/baby_janus/gateway/gateway-api"
)

func TestRedirects(t *testing.T) {
	setup()

	t.Run("Test Redirect - Simple", testRedirect)
	t.Run("Test Redirect - Cluster - Round Robin", testRedirectRoundRobin)
}

func testRedirect(t *testing.T) {

	calledBack := false
	origin := "/origin"

	http.HandleFunc("/target", func(w http.ResponseWriter, r *http.Request) {
		calledBack = true
	})

	myDomain := "http://127.0.0.1:8080"
	client := gateway_api.NewGatewayClient(myDomain)
	client.RegisterRoute(origin, fmt.Sprintf("%v/target", myDomain))

	mustHttpGet("http://127.0.0.1:8080/origin")

	if !calledBack {
		t.Error("didn't redirect to target")
	}
}

func testRedirectRoundRobin(t *testing.T) {
	t.Skip()

	var wg sync.WaitGroup

	myDomain := "http://127.0.0.1:8080"
	client := gateway_api.NewGatewayClient(myDomain)

	calledBacks := make([]bool, 100)
	for i := 0; i < len(calledBacks); i++ {
		wg.Add(1)
		go func(id int) {
			http.HandleFunc("/origin", func(w http.ResponseWriter, r *http.Request) {
				calledBacks[id] = true
			})
			client.RegisterRoute("/origin", fmt.Sprintf("%v/target%v", myDomain, id))
			wg.Done()
		}(i)
	}

	wg.Wait()

	for i := 0; i < len(calledBacks); i++ {
		mustHttpGet(fmt.Sprintf("http://127.0.0.1:8080/origin"))
	}

	for i := 0; i < len(calledBacks); i++ {
		if !calledBacks[i] {
			t.Fatal("not all callbacks were run", i)
		}
	}
}

func mustHttpPost(url string, contentType string, body []byte) []byte {
	resp, err := http.Post(url, contentType, bytes.NewBuffer(body))
	defer resp.Body.Close()
	if err != nil {
		panic(err)
	}
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return body
}

func mustHttpGet(url string) (body []byte) {
	resp, err := http.Get(url)
	defer func() {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}()
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
