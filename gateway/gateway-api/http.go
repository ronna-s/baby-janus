package gateway_api

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

func mustHttpPost(url string, contentType string, body []byte) []byte {
	resp, err := http.Post(url, contentType, bytes.NewBuffer(body))
	defer func() {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}()
	if err != nil {
		panic(fmt.Sprintf("failed calling %v reason: %v", url, err))
	}
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(fmt.Sprintf("failed reading response body. url: %v reason: %v", url, err))
	}
	return body
}


func mustHttpGet(url string) (body []byte) {
	resp, err := http.Get(url)
	defer func() {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}()
	if err != nil {
		panic(err)
	}
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return body
}
