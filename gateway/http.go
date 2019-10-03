package gateway

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

func httpPost(url string, contentType string, body []byte) error {
	resp, err := http.Post(url, contentType, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	defer resp.Body.Close()
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
