package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"reflect"
	"sync"
	"sync/atomic"

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

	// reverse proxy - create an endpoint to register new endpoints from orig to dest:
	// 1. Create an http handler for /register-endpoint path.
	// Upon request to register an endpoint add a new http handler for the orig endpoint to
	// reverse proxy the dest endpoint.
	// 2. Make the request (keep the method, headers, body and cookies of the original request)
	// 3. Don't forget to close the body
	// 4. Copy the new response Content-Type header to the response (you might want to copy more headers for caching,
	// or even all of them depending on what your apps actually do - for instance full path pagination headers may
	// not work as your endpoints are behind an API gateway).
	// 5. Copy the status code to the original response
	// 6. Copy the response body (io.Copy is your friend)
	// 7. This should work

	//[your code here!!!!]
	destinations := make(map[string][]string)
	positions := make(map[string]*int32)
	var mu sync.RWMutex
	http.HandleFunc("/register-endpoint", func(w http.ResponseWriter, r *http.Request) {
		var ep Endpoint
		if err := json.NewDecoder(r.Body).Decode(&ep); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		mu.Lock()
		defer mu.Unlock()
		_, exists := destinations[ep.Orig]
		destinations[ep.Orig] = append(destinations[ep.Orig], ep.Dest)
		var i int32 = -1
		positions[ep.Orig] = &i


		if !exists {
			http.HandleFunc(ep.Orig, func(w http.ResponseWriter, r *http.Request) {
				mu.RLock()
				defer mu.RUnlock()
				dests := destinations[ep.Orig]
				pos := atomic.AddInt32(positions[ep.Orig], 1)
				dest := dests[int(pos)%len(dests)]
				var client http.Client
				req, err := http.NewRequest(r.Method, dest, r.Body)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				//copy headers (incl. cookies)
				for k, values := range r.Header {
					for _, v := range values {
						req.Header.Add(k, v)
					}
				}

				resp, err := client.Do(req)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				defer resp.Body.Close()

				w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
				w.WriteHeader(resp.StatusCode)
				io.Copy(w, resp.Body)
			})
		}

		//your next steps go here
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
