package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"

	"github.com/wwgberlin/baby_janus/gateway/gateway-api"
	"github.com/wwgberlin/baby_janus/server/cluster"
)

func main() {
	myDomain := os.Getenv("HOSTNAME")
	client := gateway_api.NewGatewayClient("http://baby_janus_gateway:8080")
	rand.Seed(client.GetSeed())
	id := client.GetID()
	registerRoutes(client, id, myDomain)
	http.ListenAndServe(":8080", nil)
}

func registerRoutes(api gateway_api.Client, id int, myDomain string) {
	routes := getRoutes(id)
	for i := range routes {
		route := routes[i]
		http.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
			b, err := ioutil.ReadFile(fmt.Sprintf(".%s", route))
			if err != nil {
				panic(err.Error())
			}
			fmt.Fprintf(w, string(b))
		})
		api.RegisterRoute(fmt.Sprintf(route), fmt.Sprintf("http://%s:8080%s", myDomain, route))
	}
}

func getRoutes(id int) (parts []string) {
	c := cluster.NewCluster()
	if os.Getenv("CLUSTER_STRATEGY") == "all"{
		parts = c.GetParts()
	} else {
		parts = c.GetInstanceParts(id)
	}
	return
}