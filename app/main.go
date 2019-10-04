package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/ronna-s/baby-janus/gateway"
)

const numParts = 454

type (
	Endpoint struct {
		Orig string
		Dest string
	}
)

// parts - handler to get all parts
func parts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	for _, path := range getPartsURLs() {
		func() {
			route := fmt.Sprintf("http://baby-janus_gateway:8080%s", path)
			resp, err := http.Get(route)
			if err != nil {
				panic(fmt.Sprintf("error fetching %v: %v", route, err))
			}
			defer resp.Body.Close()

			// this is bad because we will have a partial response if it fails, but oh so good,
			// because it animates the response so you can visualize the process
			if _, err := io.Copy(w, resp.Body); err != nil {
				panic(err)
			}
		}()
	}

}

func getPartsURLs() []string {
	res := make([]string, numParts)
	for i := range res {
		res[i] = fmt.Sprintf("/parts/%d.part", i)
	}
	return res
}

// hello world handler
func helloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "hello world")
}

func main() {
	//myDomain := "127.0.0.1:8081"
	myDomain := os.Getenv("HOSTNAME") + ":8080"

	client := gateway.NewClient("http://baby-janus_gateway:8080")
	log.Println("registering", "/hi")
	client.RegisterRoute("/hi", fmt.Sprintf("http://%v/hi", myDomain))
	log.Println("registering", "/parts")
	client.RegisterRoute("/parts", fmt.Sprintf("http://%v/parts", myDomain))

	http.HandleFunc("/hi", helloWorld)
	http.HandleFunc("/parts", parts)

	http.ListenAndServe(":8080", nil)
}
