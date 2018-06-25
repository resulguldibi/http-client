package entity

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
)

type IHttpClient interface {
	Get(url string) IHttpClient
	PostJson(url string, data interface{}) IHttpClient
	PostUrlEncoded(url string, data url.Values) IHttpClient
	EndStruct(response interface{}) error
}

type HttpClient struct {
	client   *http.Client
	response *http.Response
}

func (client *HttpClient) Get(url string) IHttpClient {
	res, err := client.client.Get(url)
	if err != nil {
		panic(err)
	}
	client.response = res
	return client
}

func (client *HttpClient) PostJson(url string, data interface{}) IHttpClient {

	dataBytes, err := json.Marshal(data)

	if err != nil {
		panic(err)
	}

	reader := strings.NewReader(string(dataBytes))

	res, err := client.client.Post(url, "application/json", reader)

	if err != nil {
		panic(err)
	}

	client.response = res
	return client
}

func (client *HttpClient) PostUrlEncoded(url string, data url.Values) IHttpClient {

	reader := strings.NewReader(string(data.Encode()))

	res, err := client.client.Post(url, "application/x-www-form-urlencoded", reader)

	if err != nil {
		panic(err)
	}

	client.response = res
	return client
}

func (client *HttpClient) EndStruct(response interface{}) error {

	defer client.response.Body.Close()

	err := json.NewDecoder(client.response.Body).Decode(response)

	if err != nil {
		panic(err)
	}

	return err
}

func NewHttpClient() IHttpClient {
	return &HttpClient{client: &http.Client{}}
}
