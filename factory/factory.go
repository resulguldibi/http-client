package factory

import (
	"github.com/resulguldibi/http-client/entity"
)

type IHttpClientFactory interface {
	GetHttpClient() entity.IHttpClient
}

type HttpClientFactory struct {
}

func NewHttpClientFactory() IHttpClientFactory {
	return &HttpClientFactory{}
}

func (f *HttpClientFactory) GetHttpClient() entity.IHttpClient {
	return entity.NewHttpClient()
}
