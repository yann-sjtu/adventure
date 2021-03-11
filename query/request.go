package query

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"time"
)

const (
	success = "sucess"
	fail    = "failed"
)

type Request struct {
	Version string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	ID      int         `json:"id"`
}

type RPCError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type Response struct {
	Error  *RPCError       `json:"error"`
	ID     int             `json:"id"`
	Result json.RawMessage `json:"result,omitempty"`
}

func CreateRequest(method string, params interface{}) Request {
	return Request{
		Version: "2.0",
		Method:  method,
		Params:  params,
		ID:      1,
	}
}

func CallWithProxy(postBody []byte, reqType int,proxyIP string) (*Response, error) {
	client := &http.Client{}

	//是否使用代理IP
	if proxyIP != "" {
		proxy := func(_ *http.Request) (*url.URL, error) {
			return url.Parse(proxyIP)
		}
		transport := &http.Transport{Proxy: proxy}
		client = &http.Client{Transport: transport}
	} else {
		client = &http.Client{}
	}

	//post请求
	req, err := http.NewRequest("POST", host, bytes.NewBuffer(postBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	startTime := time.Now()
	resp, reqErr := client.Do(req)
	elapsed := time.Since(startTime)
	if reqErr != nil {
		log.Println(reqType, elapsed, fail, reqErr)
		return nil, reqErr
	} else {
		log.Println(reqType, elapsed, success)
	}
	//返回内容
	//var rpcRes *Response
	//decoder := json.NewDecoder(resp.Body)
	//rpcRes = new(Response)
	//err = decoder.Decode(&rpcRes)
	//if err != nil {
	//	panic(err)
	//}
	err = resp.Body.Close()
	if err != nil {
		return nil, err
	}

	//return rpcRes, nil
	return &Response{}, nil
}

func Call(request Request) (*Response, error) {
	req, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	var rpcRes *Response
	resp, err := http.Post(host, "application/json", bytes.NewBuffer(req)) //nolint:gosec
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(resp.Body)
	rpcRes = new(Response)
	err = decoder.Decode(&rpcRes)
	if err != nil {
		return nil, err
	}

	err = resp.Body.Close()
	if err != nil {
		return nil, err
	}

	return rpcRes, nil
}
