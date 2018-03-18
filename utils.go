package gin_unit_test

import (
	"net/http/httptest"
	"io/ioutil"
	"bytes"
	"github.com/gin-gonic/gin"
	"encoding/json"
	"net/http"
	"fmt"
	"os"
	"mime/multipart"
	"io"
)

var (
	// method type
	// 请求类型
	POST = "POST"
	GET = "GET"
	PUT = "PUT"
	DELETE = "DELETE"

	// the way to send parameters
	// 传参类型
	Json = "json"
	Form = "form"

	// router
	// 路由
	router *gin.Engine

	// customed request headers for token authorization and so on
	// 自定义请求头，可用于token验证等
	customedHeaders map[string]string
)

// Param interface defines two function
// Param接口定义了两个方法
type Param interface {
	// transform the request params to query string
	// 将请求参数转化为query string
	QueryStr() string

	// transform the request params to json bytes
	// 将请求参数转化为json bytes
	JsonBytes() []byte
}

// example param struct
// 示例param结构体
type EgParam struct {
	// contains request params
	// 保存请求参数的map
	Content map[string]interface{}
}

// set a param
// 设置参数
func (p *EgParam) Set(key string, value interface{}){
	p.Content[key] = value
}

// get a param
// 获取参数
func (p *EgParam) Get(key string) interface{} {
	return p.Content[key]
}

// implements the QueryStr function
// 实现了QueryStr方法
func (p *EgParam) QueryStr() string {
	if p == nil {
		return ""
	}
	values := ""
	for key, val := range p.Content {
		values += "&" + key + "=" + fmt.Sprintf("%v", val)
	}
	return values[1:]
}

// implements the JsonBytes function
// 实现了JsonBytes方法
func (p *EgParam) JsonBytes() []byte {
	jsonBytes,_ := json.Marshal(p.Content)
	return jsonBytes
}

// SetRouter
// 设置路由
func SetRouter(r *gin.Engine) {
	router = r
}

// SetCustomedHeaders
// 设置自定义请求头
func SetCustomedHeaders(h map[string]string) {
	customedHeaders = h
}

// simulate sending the request and accept the response, return the response body
// 模拟发起请求，并接收响应，返回响应中的body
func runHandler(req *http.Request) (bodyByte []byte, err error) {

	// initialize response record
	// 初始化response record
	w := httptest.NewRecorder()

	// simulate running the handler
	// 模拟请求的发送，调用相应的handler处理请求
	router.ServeHTTP(w, req)

	// extract the response from the response record
	// 从response record中提取response
	result := w.Result()
	defer result.Body.Close()

	// extract response body
	// 从response中提取response body
	bodyByte,err = ioutil.ReadAll(result.Body)

	return
}

// simulate sending file and receive the response, extract the response body
// the first parameter is the method, such as POST,PUT
// the second parameter is the request uri
// the third parameter is the name of the file, containing the directory of the file
// the forth parameter is the field name of the file
// 模拟文件的发送和响应的接收，并从响应中提取响应的body部分
// 第一个参数是请求的方法，是POST，PUT中的一种
// 第二个参数是请求地址
// 第三个参数是文件名
// 第四个参数是文件对应的字段名
func SimulateFileHandler(method,uri,fileName string, filedName string) (bodyByte []byte, err error){
	// check whether the router is nil
	// 检查router是否为空
	if router == nil {
		err = ErrRouterNotSet
		return
	}

	// check whether the method is appropiate, now the method must should be POST or PUT
	// 检查method是否合适，此时的method应该为POST或PUT中的一种
	if method != POST && method != PUT {
		err = ErrMustPostOrPut
		return
	}

	// create form file
	// 创建表单文件
	buf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(buf)

	fileWriter,err := bodyWriter.CreateFormFile(filedName, fileName)
	if err != nil {
		return
	}

	// open the file
	// 打开文件
	file,err := os.Open(fileName)
	if err != nil {
		return
	}

	// copy the content of the file to the fileWriter
	// 将文件的内容拷贝到fileWriter中，其实也就相当于拷贝到了之前的buf中了
	length,err := io.Copy(fileWriter, file)
	if err != nil {
		return
	}

	file.Close()
	bodyWriter.Close()

	// make request
	// 构造请求
	req := httptest.NewRequest(method, uri, buf)
	req.Header.Add("Content-Type", bodyWriter.FormDataContentType())
	err = req.ParseMultipartForm(length)

	if err != nil {
		return
	}

	// check whether the customed headers are nil
	// if not, then add them to current request headers
	// 判断自定义请求头是否为空，如果不为空，则将他们添加到现有的请求头中
	if customedHeaders != nil {
		// add the customed headers
		// 添加自定义头信息
		for key,value := range customedHeaders {
			req.Header.Set(key, value)
		}
	}

	// simulate sending request
	// 模拟发送请求
	bodyByte, err = runHandler(req)
	return
}

// simulate sending request and receive the response, extract the response body
// the first parameter is the method, such as GET,POST,PUT,DELETE
// the second parameter is the request uri
// the third parameter is the way to send parameter, expected form or json
// the forth parameter is the parameter of the request, it must implement the Param interface
// 模拟请求的发送和响应的接收，并从响应中提取响应的body部分
// 第一个参数是请求的方法，是GET，POST，PUT，DELETE中的一种
// 第二个参数是请求地址
// 第三个参数是发送请求参数的方式，是json或form表单中的一种
// 第四个参数是请求参数，必须实现了Param接口
func SimulateOrdinaryHandler(method string, uri string, way string, param Param) (bodyByte []byte, err error){
	if router == nil {
		err = ErrRouterNotSet
		return
	}

	// make request
	// 构造请求
	var req *http.Request
	switch way {
	case Json:
		// when the way is Json, change the variable to json bytes, and add to the request body
		// 当传参方式为json时，将变量转化为json bytes，添加到请求体中
		jsonBytes := []byte{}
		if param != nil {
			jsonBytes = param.JsonBytes()
		}
		req = httptest.NewRequest(method, uri, bytes.NewReader(jsonBytes))
	case Form:
		// when the way is form, then change the variable to query string, and add to the latter of the request uri
		// 当传参方式为form时，将变量转化为query string	，添加到请求uri之后
		queryStr := ""
		if param != nil {
			queryStr = param.QueryStr()
		}
		req = httptest.NewRequest(method, uri+"?"+queryStr, nil)
	}

	// check whether the customed headers are nil
	// if not, then add them to current request headers
	// 判断自定义请求头是否为空，如果不为空，则将他们添加到现有的请求头中
	if customedHeaders != nil {
		// add the customed headers
		// 将自定义头添加到请求头中
		for key,value := range customedHeaders {
			req.Header.Set(key, value)
		}
	}

	// launch request
	// 发起请求
	bodyByte, err = runHandler(req)
	return
}