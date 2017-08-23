package gateway_api

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

func mustHttpPost(url string, contentType string, body []byte) (err error) {
	var resp *http.Response
	resp, err = http.Post(url, contentType, bytes.NewBuffer(body))
	defer func() {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}()
	if err != nil {
		return
	}
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	return nil
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
