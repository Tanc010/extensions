package main

import (
	"context"

	"github.com/mosn/extensions/go-plugin/pkg/protocol/bolt"
	"github.com/valyala/fasthttp"
	"mosn.io/api"
	at "mosn.io/api/extensions/transcoder"
	"mosn.io/pkg/protocol/http"
)

type bolt2sp struct {
	cfg map[string]interface{}
}

func NewTranscoder(cfg map[string]interface{}) at.Transcoder {
	return &bolt2sp{
		cfg: cfg,
	}
}

func (t *bolt2sp) Accept(ctx context.Context, headers api.HeaderMap, buf api.IoBuffer, trailers api.HeaderMap) bool {
	return true
}

func (t *bolt2sp) TranscodingRequest(ctx context.Context, headers api.HeaderMap, buf api.IoBuffer, trailers api.HeaderMap) (api.HeaderMap, api.IoBuffer, api.HeaderMap, error) {
	sourceRequest, ok := headers.(*bolt.Request)
	if !ok {
		return headers, buf, trailers, nil
	}
	targetRequest := fasthttp.Request{}
	// 1. headers
	sourceRequest.Range(func(Key, Value string) bool {
		targetRequest.Header.Set(Key, Value)
		return true
	})
	// 协议头变更
	path := t.cfg["x-mosn-path"].(string)
	targetRequest.Header.Set("x-mosn-path", path)
	methond := t.cfg["x-mosn-method"].(string)
	targetRequest.Header.Set("x-mosn-methond", methond)
	return http.RequestHeader{RequestHeader: &targetRequest.Header}, buf, trailers, nil
}

func (t *bolt2sp) TranscodingResponse(ctx context.Context, headers api.HeaderMap, buf api.IoBuffer, trailers api.HeaderMap) (api.HeaderMap, api.IoBuffer, api.HeaderMap, error) {
	sourceResponse, ok := headers.(http.ResponseHeader)
	if !ok {
		return headers, buf, trailers, nil
	}
	targetResponse := bolt.NewRpcResponse(0, uint16(sourceResponse.StatusCode()), headers, buf)
	return targetResponse, buf, trailers, nil

}
