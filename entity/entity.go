package entity

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
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
	err      error
	inTrx    bool
}

func (client *HttpClient) Get(url string) IHttpClient {

	if client.inTrx {
		panic(errors.New("client is in trx"))
	}

	client.inTrx = true

	res, err := client.client.Get(url)
	if err != nil {
		client.err = errors.New(fmt.Sprintf("error in Get-> %s", err.Error()))
		return client
	}

	client.response = res
	return client
}

func (client *HttpClient) PostJson(url string, data interface{}) IHttpClient {

	if client.inTrx {
		panic(errors.New("client is in trx"))
	}

	client.inTrx = true

	dataBytes, err := json.Marshal(data)

	if err != nil {
		client.err = errors.New(fmt.Sprintf("error in PostJson marshal-> %s", err.Error()))
		return client
	}

	reader := strings.NewReader(string(dataBytes))

	res, err := client.client.Post(url, "application/json", reader)

	if err != nil {
		client.err = errors.New(fmt.Sprintf("error in Post-> %s", err.Error()))
		return client
	}

	client.response = res
	return client
}

func (client *HttpClient) PostUrlEncoded(url string, data url.Values) IHttpClient {

	if client.inTrx {
		panic(errors.New("client is in trx"))
	}

	client.inTrx = true

	reader := strings.NewReader(string(data.Encode()))

	res, err := client.client.Post(url, "application/x-www-form-urlencoded", reader)

	if err != nil {
		client.err = errors.New(fmt.Sprintf("error in response decoding-> %s", err.Error()))
		return client
	}

	client.response = res
	return client
}

func (client *HttpClient) EndStruct(response interface{}) error {

	if !client.inTrx {
		panic(errors.New("client is not in trx"))
	}

	client.inTrx = false

	if client.err != nil {
		return client.err
	}

	client.err = nil

	if client.response.StatusCode != http.StatusOK {
		defer client.response.Body.Close()
		resp, err := ioutil.ReadAll(client.response.Body)

		if err != nil {
			fmt.Println("ioutil.ReadAll error ->", err.Error())
		} else {
			fmt.Println("client.response.Body ->", string(resp))
		}

		return errors.New(fmt.Sprintf("response.StatusCode is not OK -> %d", client.response.StatusCode))
	}

	defer client.response.Body.Close()

	err := json.NewDecoder(client.response.Body).Decode(response)

	if err != nil {
		return errors.New(fmt.Sprintf("error in response decoding-> %s", err.Error()))
	}

	return err
}

func NewHttpClient() IHttpClient {
	return &HttpClient{client: &http.Client{}}
}
