package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"

	"github.com/wwgberlin/baby_janus/gateway/cluster"
)

type (
	Endpoint struct {
		Origin string
		Target string
	}
	api struct {
		routes     map[string]chan string
		incomingCh chan Endpoint
	}
)

func main() {
	a := api{routes: map[string]chan string{}, incomingCh: make(chan Endpoint, 1)}
	c := cluster.NewCluster()

	//gateway api operations
	http.HandleFunc("/routes", listRoutes)
	http.HandleFunc("/register_endpoint", a.registerEndpoint)

	//cluster operations
	http.HandleFunc("/next_cluster_id", incrClusterId(c))
	http.HandleFunc("/seed", getSeed(c))

	http.ListenAndServe(":8080", nil)
}

/*
	registerEndpoint handles requests to registers routes origin - target
 */
func (a *api) registerEndpoint(response http.ResponseWriter, request *http.Request) {
	var endpoint Endpoint
	if request != nil && request.Body != nil {
		body, err := ioutil.ReadAll(request.Body)
		if err != nil {
			response.WriteHeader(http.StatusBadRequest)
			return
		}
		if err := json.Unmarshal(body, &endpoint); err != nil {
			response.WriteHeader(http.StatusBadRequest)
		}
		http.HandleFunc(endpoint.Origin, a.redirectHandler(endpoint.Target))
		response.WriteHeader(http.StatusCreated)
	}
}

/*
	redirectHandler returns a handler to redirect the request to
 */

func (a *api) redirectHandler(target string) func(http.ResponseWriter, *http.Request) {
	fmt.Println("registered", target)
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, target, http.StatusFound)
	}
}


/*
	handler to list all routes the gateway is responding to
 */
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

/*
	seed handler for cluster randomization operations
 */

func getSeed(c cluster.Cluster) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, fmt.Sprintf("%v", c.GetSeed()))
	}
}

/*
	incrClusterId - returns handler to increment the cluster servers size
 */

func incrClusterId(c cluster.Cluster) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, fmt.Sprintf("%v", c.IncrClusterId()))
	}
}
