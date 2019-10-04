package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
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
	// 1. Create an http handler for /register-endpoint. Upon request to register an endpoint add a new http handler
	// for the orig endpoint that will retrieve the dest endpoint
	// 2. Make the request (keep the method, headers, body and cookies of the original request)
	// 3. Don't forget to close the body
	// 4. Copy the new response Content-Type header to the response (you might want to copy more headers for caching,
	// or even all of them depending on what your apps actually do - for instance full path pagination headers may
	// not work as your endpoints are behind an API gateway).
	// 5. Copy the status code to the original response
	// 6. Copy the response body (io.Copy is your friend)
	// 7. This should work

	http.HandleFunc("/register-endpoint", func(w http.ResponseWriter, r *http.Request) {
		var e Endpoint
		if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		http.HandleFunc(e.Orig, func(w http.ResponseWriter, r *http.Request) {
			var client http.Client
			req, err := http.NewRequest(r.Method, e.Dest, r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			for k, values := range r.Header {
				for _, v := range values {
					req.Header.Add(k, v)
				}
			}
			for _, c := range r.Cookies() {
				req.AddCookie(c)
			}

			resp, err := client.Do(req)
			defer resp.Body.Close()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
			w.WriteHeader(resp.StatusCode)
			io.Copy(w, resp.Body)
		})
		w.WriteHeader(http.StatusCreated)

	})

	go func() {
		close(ready)
	}()
	log.Fatal(http.ListenAndServe(":8080", nil))

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
