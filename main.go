package main

import (
	"net/http"
	"bytes"
	"fmt"
	"io/ioutil"

	"context"
	"encoding/json"
	"log"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

func httpPost(url, data string) string {
    
	// url := "http://localhost:8080/mul"
    fmt.Println("URL:>", url)

    //json序列化
	// post := "{\"x\":119,\"y\":100}"
	post := data

    fmt.Println(url, "post", post)

    var jsonStr = []byte(post)
    fmt.Println("jsonStr", jsonStr)
    fmt.Println("new_str", bytes.NewBuffer(jsonStr))

    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
    // req.Header.Set("X-Custom-Header", "myvalue")
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    fmt.Println("response Status:", resp.Status)
    fmt.Println("response Headers:", resp.Header)
    body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
	return string(body)
}

type CallService interface {
	Call(context.Context, string, string) string
}

type CallServiceImpl struct {}
func (CallServiceImpl) Call(_ context.Context, url, data string) string {
	return httpPost(url, data)
}

type callRequest struct {
	URL string `json:"url"`
	DATA string `json:"data"`
}
type callResponse struct {
	S string `json:"s"`
}

// 提供给 NewServer 函数
func makeCallEndpoint(csvc CallService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(callRequest)
		S := csvc.Call(ctx, req.URL, req.DATA)
		return callResponse{S}, nil
	}
}

func decodeCallRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request callRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

// 主入口
func main() {
	csvc :=  CallServiceImpl{}

	callHandler := httptransport.NewServer(
		makeCallEndpoint(csvc),
		decodeCallRequest,
		encodeResponse,
	)

	// 路由配置
	http.Handle("/", callHandler)
	
	// 日志
	log.Fatal(http.ListenAndServe(":9090", nil))
}
