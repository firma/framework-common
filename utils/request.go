package utils

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// 请求客户端
var httpClient = &http.Client{}

func RequestGet(reqUrl string, reqParams map[string]interface{}, headers map[string]string) (respData []byte, err error) {
	urlParams := url.Values{}
	Url, _ := url.Parse(reqUrl)
	for key, val := range reqParams {
		urlParams.Set(key, fmt.Sprintf("%v", val))
	}
	//如果参数中有中文参数,这个方法会进行URLEncode
	Url.RawQuery = urlParams.Encode()
	urlPath := Url.String()

	httpRequest, _ := http.NewRequest("GET", urlPath, nil)
	// 添加请求头
	if headers != nil {
		for k, v := range headers {
			httpRequest.Header.Add(k, v)
		}
	}
	// 发送请求
	resp, err := httpClient.Do(httpRequest)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	var bodyBytes []byte
	if resp.Body != nil {
		bodyBytes, _ = io.ReadAll(resp.Body)
	}
	return bodyBytes, nil
}
