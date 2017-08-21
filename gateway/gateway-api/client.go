package gateway_api

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type (
	Client interface {
		GetID() int
		GetSeed() int64
		RegisterRoute(origin string, target string)
	}
	client struct {
		gatewayServer string
	}
)

func NewGatewayClient(gatewayServer string) Client {
	return &client{gatewayServer: gatewayServer}
}

func (client *client) RegisterRoute(origin string, target string) {
	b, err := json.Marshal(struct {
		Origin string
		Target string
	}{Origin: origin, Target: target})

	if err != nil {
		panic(fmt.Sprintf("failed registering route %s => %s reason: %s", origin, target, err.Error()))
	}

	mustHttpPost(fmt.Sprintf("%s/register_endpoint", client.gatewayServer), "application/json", b)
}

func (client *client) GetID() int {
	body := mustHttpGet(fmt.Sprintf("%s/next_cluster_id", client.gatewayServer))
	id, err := strconv.Atoi(string(body))
	if err != nil {
		panic("could not init id")
	}
	return id
}

func (client *client) GetSeed() int64 {
	body := mustHttpGet(fmt.Sprintf("%s/seed", client.gatewayServer))
	seed, err := strconv.ParseInt(string(body), 10, 64)
	if err != nil {
		panic("can't init seed")
	}
	return seed
}