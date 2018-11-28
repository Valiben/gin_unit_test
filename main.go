package gin_unit_test

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"zonst/qipai/gin-unittest-demo/utils"
)

var (
	// router
	router http.Handler

	// customed request headers for token authorization and so on
	myHeaders = make(map[string]string, 0)

	logging *log.Logger
)

// set the router
func SetRouter(r http.Handler) {
	router = r
}

// set the log
func SetLog(l *log.Logger) {
	logging = l
}

// add custom request header
func AddHeader(key, value string) {
	myHeaders[key] = value
}

// printf log
func printfLog(format string, v ...interface{}) {
	if logging == nil {
		return
	}

	logging.Printf(format, v...)
}

// invoke handler
func invokeHandler(req *http.Request) (bodyByte []byte, err error) {

	// initialize response record
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// extract the response from the response record
	result := w.Result()
	defer result.Body.Close()

	// extract response body
	bodyByte, err = ioutil.ReadAll(result.Body)

	return
}

func TestFileHandler(method, api, fileName string, fieldName string, param interface{}) (bodyByte []byte, err error) {
	// check whether the router is nil
	if router == nil {
		err = ErrRouterNotSet
		return
	}

	_, paramStr := utils.MakeQueryStrFrom(param)
	printfLog("TestFileHandler\tRequest:\t%v:%v?%v \tFileName:%v, FieldName:%v\n",
		method, api, paramStr, fileName, fieldName)

	// make request
	req, err := utils.MakeFileRequest(method, api, fileName, fieldName, param)
	if err != nil {
		return
	}

	for key, value := range myHeaders {
		req.Header.Add(key, value)
	}

	// invoke handler
	bodyByte, err = invokeHandler(req)
	printfLog("TestFileHandler\tResponse:\t%v:%v,\tResponse:%v\n\n\n", method, api, string(bodyByte))
	return
}

func TestOrdinaryHandler(method string, api string, mime string, param interface{}) (bodyByte []byte, err error) {
	if router == nil {
		err = ErrRouterNotSet
		return
	}

	printfLog("TestOrdinaryHandler\tRequest:\t%v:%v,\trequestBody:%v\n", method, api, param)

	// make request
	req, err := utils.MakeRequest(method, mime, api, param)
	if err != nil {
		return
	}

	// add the customed headers
	for key, value := range myHeaders {
		req.Header.Add(key, value)
	}

	// invoke handler
	bodyByte, err = invokeHandler(req)

	printfLog("TestOrdinaryHandler\tResponse:\t%v:%v\tResponse:%v\n\n\n", method, api, string(bodyByte))
	return
}

func TestHandlerUnMarshalResp(method string, uri string, way string, param interface{}, resp interface{}) error {
	bodyByte, err := TestOrdinaryHandler(method, uri, way, param)
	if err != nil {
		return err
	}

	return json.Unmarshal(bodyByte, resp)
}

func TestFileHandlerUnMarshalResp(method, uri, fileName string, filedName string, param interface{}, resp interface{}) error {
	bodyByte, err := TestFileHandler(method, uri, fileName, filedName, param)
	if err != nil {
		return err
	}

	return json.Unmarshal(bodyByte, resp)
}
