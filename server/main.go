package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/ronna-s/baby-janus/gateway"
)

func init(){
	time.Sleep(15*time.Second)
}

func main() {
	myDomain := os.Getenv("HOSTNAME")
	client := gateway.NewClient("http://baby-janus_gateway:8080")
	rand.Seed(client.GetSeed())
	id := client.GetID()
	registerRoutes(client, id, myDomain)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func registerRoutes(api *gateway.Client, id int, myDomain string) {
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
		fmt.Printf("registering route %s\n", route)
		api.RegisterRoute(fmt.Sprintf(route), fmt.Sprintf("http://%s:8080%s", myDomain, route))
	}
}

func getRoutes(id int) (parts []string) {
	c := NewCluster()
	if os.Getenv("CLUSTER_STRATEGY") == "all" {
		parts = c.GetParts()
	} else {
		parts = c.GetInstanceParts(id)
	}
	return
}
