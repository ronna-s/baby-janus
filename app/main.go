package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/wwgberlin/baby_janus/gateway/gateway-api"
)

type (
	Endpoint struct {
		Origin string
		Target string
	}
)

const NUM_PARTS = 136

/*
	parts - handler to get all parts
 */
func parts(w http.ResponseWriter, r *http.Request) {
	bodies := []string{}
	for _, path := range getPartsURLs() {
		route := fmt.Sprintf("http://baby_janus_gateway:8080%s", path)
		resp, err := http.Get(route)
		if err != nil {
			panic(fmt.Sprintf("error fetching %v: %v", route, err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		bodies = append(bodies, string(body))
	}
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, strings.Join(bodies, ""))
}

func getPartsURLs() []string {
	res := make([]string, NUM_PARTS)
	for i := range res {
		res[i] = fmt.Sprintf("/parts/%d.part", i)
	}
	return res
}

/*
	hello world handler
 */
func helloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello world")
}

func main() {
	myDomain := os.Getenv("VIRTUAL_HOST")
	if myDomain == "" {
		myDomain = "127.0.0.1:8081"
	}
	client := gateway_api.NewGatewayClient("http://baby_janus_gateway:8080")
	fmt.Println("asking to register", "/hi")
	client.RegisterRoute("/hi", fmt.Sprintf("http://%v/hi", myDomain))
	fmt.Println("asking to register", "/parts")
	client.RegisterRoute("/parts", fmt.Sprintf("http://%v/parts", myDomain))

	http.HandleFunc("/hi", helloWorld)
	http.HandleFunc("/parts", parts)

	http.ListenAndServe(":8080", nil)
}
