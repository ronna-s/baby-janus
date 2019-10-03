package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/ronna-s/baby-janus/gateway"
)

type (
	Endpoint struct {
		Orig string
		Dest string
	}
)

func main() {
	ready := make(chan struct{})
	run(ready)
}

func run(ready chan struct{}) {
	c := gateway.NewCluster()

	//gateway api operations
	http.HandleFunc("/routes", listRoutes)

	//cluster operations
	http.HandleFunc("/next_cluster_id", incrClusterId(c))
	http.HandleFunc("/seed", getSeed(c))

	// reverse proxy - register a new endpoint from orig to dest:
	// 1. Create an http request to the destination endpoint using the method of the a original request
	// 		you will need a client and to copy well important stuff like cookies to the new request
	// 2. Make the request
	// 3. Don't forget to close the body
	// 4. Copy the new response Content-Type header to the response (you might want to copy more headers for caching,
	// or even all of them depending on what your apps actually do - for instance full path pagination headers may
	// not work as your endpoints are behind an API gateway).
	// 5. Copy the status code to the original response
	// 6. Copy the response body
	// 7. This should work

	go func() {
		close(ready)
	}()
	http.ListenAndServe(":8080", nil)

}

// handler to list all routes the gateway is responding to
func listRoutes(w http.ResponseWriter, r *http.Request) {
	httpMux := reflect.ValueOf(http.DefaultServeMux).Elem()
	finList := httpMux.FieldByIndex([]int{1})
	keys := finList.MapKeys()
	routes := make([]string, len(keys))
	for i := range keys {
		routes[i] = fmt.Sprintf("%v", keys[i])
	}
	b, _ := json.Marshal(routes)
	fmt.Fprintf(w, string(b))
}

// seed handler for cluster randomization operations
func getSeed(c gateway.Cluster) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, fmt.Sprintf("%v", c.GetSeed()))
	}
}

// incrClusterId - returns handler to increment the cluster servers size
func incrClusterId(c gateway.Cluster) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, fmt.Sprintf("%v", c.IncrClusterId()))
	}
}
