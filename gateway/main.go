package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/wwgberlin/baby_janus/gateway/cluster"
)

type (
	Endpoint struct {
		Orig string
		Dest string
	}
)

func main() {
	c := cluster.NewCluster()

	//gateway api operations
	http.HandleFunc("/routes", listRoutes)

	//cluster operations
	http.HandleFunc("/next_cluster_id", incrClusterId(c))
	http.HandleFunc("/seed", getSeed(c))

	http.ListenAndServe(":8080", nil)
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
