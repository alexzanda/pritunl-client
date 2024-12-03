package pritunl

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// Config pritunl客户端配置
type Config struct {
	ApiToken     string // api token
	ApiSecret    string // api secret
	HttpProtocol string // http协议类型，默认是https
	Host         string // pritunl主机地址
	Context      *context.Context
}

// Client pritunl客户端
type Client struct {
	config     Config
	endpoint   string // http url base path
	context    context.Context
	httpClient *http.Client
}

// NewClient 获取pritunl客户端
func NewClient(apiToken, apiSecret, host string, context context.Context) (*Client, error) {
	if len(apiToken) == 0 {
		return nil, errors.New("api token不能为空")
	}
	if len(apiSecret) == 0 {
		return nil, errors.New("api secret不能为空")
	}

	httpClient := http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}}

	client := Client{
		config: Config{
			ApiToken:  apiToken,
			ApiSecret: apiSecret,
			Host:      host,
		},
		endpoint:   fmt.Sprintf("https://%s", host),
		httpClient: &httpClient,
	}
	if context != nil {
		client.context = context
	}
	return &client, nil
}

// serverUrl 返回完整的请求url
func (c *Client) serverUrl(parts ...string) string {
	return c.endpoint + strings.Join(parts, "/")
}

var applicationJSON = "application/json"

type RequestOpts struct {
	// RequestParam 查询参数
	RequestParam map[string]string
	// JSONBody 可为空，如果不为空的话，content-type变为application/json，不能和RawBody同时使用
	JSONBody interface{}
	// RawBody 可为空，不为空的话，会直接赋值给http Request，且不会额外加application/json的头
	RawBody io.Reader
	// JSONResponse 如果被指定的话，响应体将会被解析到这个字段
	JSONResponse interface{}
	// OkCodes 指定属于正常响应的状态码
	OkCodes []int
	// MoreHeaders 特定的http 头，如果不指定，会被默认的请求头覆盖
	MoreHeaders map[string]string
	// OmitHeaders 默认的请求头
	OmitHeaders []string
	// KeepResponseBody 是否保留原始的响应体，如果上层业务有进一步解析的需求的话，应该置为true
	KeepResponseBody bool
}

// Request 执行具体的请求
func (c *Client) Request(method, path string, options *RequestOpts) (*http.Response, error) {
	// 添加认证头
	authHeader, err := generateAuthHeader(path, method, c.config.ApiToken, c.config.ApiSecret)
	if err != nil {
		return nil, err
	}
	if options == nil {
		options = &RequestOpts{}
	}
	options.MoreHeaders = authHeader
	return c.doRequest(strings.ToUpper(method), c.serverUrl(path), options)
}

// doRequest 真正执行请求
func (c *Client) doRequest(method, fullUrl string, options *RequestOpts) (*http.Response, error) {
	var body io.Reader
	var contentType *string

	// 处理请求体，如果有json数据，默认的内容类型会改为json
	if options.JSONBody != nil {
		if options.RawBody != nil {
			return nil, errors.New("please provide only one of JSONBody or RawBody")
		}

		rendered, err := json.Marshal(options.JSONBody)
		if err != nil {
			return nil, err
		}

		body = bytes.NewReader(rendered)
		contentType = &applicationJSON
	}

	// KeepResponseBody和JSONResponse不能共用
	if options.KeepResponseBody && options.JSONResponse != nil {
		return nil, errors.New("cannot use KeepResponseBody when JSONResponse is not nil")
	}

	if options.RawBody != nil {
		body = options.RawBody
	}

	// 如果有查询参数，那么拼接到url中
	if options.RequestParam != nil {
		param := url.Values{}
		for key, val := range options.RequestParam {
			param.Add(key, val)
		}
		fullUrl += fmt.Sprintf("?%s", param.Encode())
	}

	// 构造http请求
	req, err := http.NewRequest(method, fullUrl, body)
	if err != nil {
		return nil, err
	}
	if c.context != nil {
		req = req.WithContext(c.context)
	}

	// 设置内容头
	if contentType != nil {
		req.Header.Set("Content-Type", *contentType)
	}
	req.Header.Set("Accept", applicationJSON)

	// 如果有更多的头就拼接上去
	if options.MoreHeaders != nil {
		for k, v := range options.MoreHeaders {
			req.Header.Set(k, v)
		}
	}

	for _, v := range options.OmitHeaders {
		req.Header.Del(v)
	}

	// 发出请求
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	// 如果没有显示指定ok的状态吗，那么就使用默认的状态吗
	okc := options.OkCodes
	if okc == nil {
		okc = defaultOkCodes(method)
	}

	// 验证http响应状态
	var ok bool
	for _, code := range okc {
		if resp.StatusCode == code {
			ok = true
			break
		}
	}
	if !ok {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("read resp body failed, err: %s", err.Error())
		}
		return resp, fmt.Errorf("resp status code is: %d, resp body: %s", resp.StatusCode, string(body))
	}

	// 如果需要的话，解析json响应体
	if options.JSONResponse != nil {
		defer resp.Body.Close()
		// 如果响应体没内容就不解析
		if resp.StatusCode == http.StatusNoContent {
			// 把响应体读干净，否则http连接将会被关闭并且无法被使用
			_, err = io.Copy(io.Discard, resp.Body)
			return resp, err
		}
		if err = json.NewDecoder(resp.Body).Decode(options.JSONResponse); err != nil {
			return nil, err
		}
	}

	// 关闭未使用的http body，以便http连接被复用
	if !options.KeepResponseBody && options.JSONResponse == nil {
		defer resp.Body.Close()
		if _, err = io.Copy(io.Discard, resp.Body); err != nil {
			return nil, err
		}
	}

	return resp, nil
}

func defaultOkCodes(method string) []int {
	switch method {
	case "GET", "HEAD":
		return []int{200}
	case "POST":
		return []int{200, 201, 202}
	case "PUT":
		return []int{200, 201, 202}
	case "PATCH":
		return []int{200, 202, 204}
	case "DELETE":
		return []int{200, 202, 204}
	}

	return []int{}
}
