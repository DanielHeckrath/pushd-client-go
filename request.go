package client

import (
	"io"
	"io/ioutil"
	"net/http"
)

func (this *Request) get(path string) (int, string, error) {
	url := this.endpoint + path
	res, err := this.client.Get(url)
	if err != nil {
		return 0, "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return 0, "", err
	}

	return res.StatusCode, string(body), nil
}

func (this *Request) post(path, contentType string, payload io.Reader) (int, string, error) {
	url := this.endpoint + path
	res, err := this.client.Post(url, contentType, payload)
	if err != nil {
		return 0, "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return 0, "", err
	}

	return res.StatusCode, string(body), nil
}

func (this *Request) del(path string) (int, string, error) {
	url := this.endpoint + path
	req, err := http.NewRequest("DELETE", url, nil)

	if err != nil {
		return 0, "", err
	}

	res, err := this.client.Do(req)
	if err != nil {
		return 0, "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return 0, "", err
	}

	return res.StatusCode, string(body), nil
}
