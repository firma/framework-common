package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

type UploadFile struct {
	// 表单名称
	Name string
	// 文件全路径
	Filepath string
}


func HttpGet(reqUrl string, reqParams map[string]string, headers map[string]string) ([]byte, error) {
	urlParams := url.Values{}
	Url, _ := url.Parse(reqUrl)
	for key, val := range reqParams {
		urlParams.Set(key, val)
	}

	//如果参数中有中文参数,这个方法会进行URLEncode
	Url.RawQuery = urlParams.Encode()
	// 得到完整的url，http://xx?query
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
		return nil, err
	}
	defer resp.Body.Close()
	response, err := io.ReadAll(resp.Body)
	return response, err
}

func PostForm(reqUrl string, reqParams map[string]interface{}, headers map[string]string) ([]byte, error) {
	return post(reqUrl, reqParams, "application/x-www-form-urlencoded", nil, headers)
}

func PostJson(reqUrl string, reqParams map[string]interface{}, headers map[string]string) ([]byte, error) {
	return post(reqUrl, reqParams, "application/json", nil, headers)
}

func PostFile(reqUrl string, reqParams map[string]interface{}, files []UploadFile, headers map[string]string) ([]byte, error) {
	return post(reqUrl, reqParams, "multipart/form-data", files, headers)
}

func DelJson(reqUrl string, body map[string]interface{}, headers map[string]string) (respData []byte, err error) {
	return del(reqUrl, body, "multipart/form-data", headers)
}

func del(reqUrl string, reqParams map[string]interface{}, contentType string, headers map[string]string) ([]byte, error) {
	requestBody, realContentType := getDelReader(reqParams, contentType)
	httpRequest, _ := http.NewRequest("POST", reqUrl, requestBody)
	// 添加请求头
	httpRequest.Header.Add("Content-Type", realContentType)
	if headers != nil {
		for k, v := range headers {
			httpRequest.Header.Add(k, v)
		}
	}
	// 发送请求
	resp, err := httpClient.Do(httpRequest)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func post(reqUrl string, reqParams map[string]interface{}, contentType string, files []UploadFile, headers map[string]string) ([]byte, error) {
	requestBody, realContentType := getReader(reqParams, contentType, files)
	httpRequest, _ := http.NewRequest("DELETE", reqUrl, requestBody)
	// 添加请求头
	httpRequest.Header.Add("Content-Type", realContentType)
	if headers != nil {
		for k, v := range headers {
			httpRequest.Header.Add(k, v)
		}
	}
	// 发送请求
	resp, err := httpClient.Do(httpRequest)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func PostByte(reqUrl string, reqParams map[string]string, files []byte, headers map[string]string) ([]byte, error) {
	//requestBody, realContentType := getReader(reqParams, contentType, files)
	body := &bytes.Buffer{}
	// 文件写入 body
	writer := multipart.NewWriter(body)
	reader := bytes.NewBuffer(files)
	part, err := writer.CreateFormFile("file", filepath.Base("file.mp3"))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, reader)

	// 其他参数列表写入 body
	if len(reqParams) > 0 {
		for k, v := range reqParams {
			if err := writer.WriteField(k, v); err != nil {
				return nil, err
			}
		}

	}
	if err := writer.Close(); err != nil {

	}
	// 上传文件需要自己专用的contentType
	realContentType := writer.FormDataContentType()
	httpRequest, _ := http.NewRequest("POST", reqUrl, body)
	// 添加请求头
	httpRequest.Header.Add("Content-Type", realContentType)
	if headers != nil {
		for k, v := range headers {
			httpRequest.Header.Add(k, v)
		}
	}
	// 发送请求
	resp, err := httpClient.Do(httpRequest)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
	//return string(response)
}

func getDelReader(reqParams map[string]interface{}, contentType string) (io.Reader, string) {
	if strings.Index(contentType, "json") > -1 {
		bytesData, _ := json.Marshal(reqParams)
		return bytes.NewReader(bytesData), contentType
	} else {
		urlValues := url.Values{}
		for key, val := range reqParams {
			urlValues.Set(key, fmt.Sprintf("%v", val))
		}
		reqBody := urlValues.Encode()
		return strings.NewReader(reqBody), contentType
	}
}

func getReader(reqParams map[string]interface{}, contentType string, files []UploadFile) (io.Reader, string) {
	if strings.Index(contentType, "json") > -1 {
		bytesData, _ := json.Marshal(reqParams)
		return bytes.NewReader(bytesData), contentType
	} else if files != nil {
		body := &bytes.Buffer{}
		// 文件写入 body
		writer := multipart.NewWriter(body)
		for _, uploadFile := range files {
			file, err := os.Open(uploadFile.Filepath)
			if err != nil {
				return nil, ""
			}
			part, err := writer.CreateFormFile(uploadFile.Name, filepath.Base(uploadFile.Filepath))
			if err != nil {
				return nil, ""
			}
			_, err = io.Copy(part, file)
			file.Close()
		}
		// 其他参数列表写入 body
		for k, v := range reqParams {
			if err := writer.WriteField(k, fmt.Sprintf("%v", v)); err != nil {
				return nil, ""
			}
		}
		if err := writer.Close(); err != nil {
			return nil, ""
		}
		// 上传文件需要自己专用的contentType
		return body, writer.FormDataContentType()
	} else {
		urlValues := url.Values{}
		for key, val := range reqParams {
			urlValues.Set(key, fmt.Sprintf("%v", val))
		}
		reqBody := urlValues.Encode()
		return strings.NewReader(reqBody), contentType
	}
}
