package gin_unit_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
)

var (
	// router
	router *gin.Engine

	// customed request headers for token authorization and so on
	myHeaders = []CustomizedHeader{}

	logging *log.Logger
)

// set the router
func SetRouter(r *gin.Engine) {
	router = r
}

// set the log
func SetLog(l *log.Logger) {
	logging = l
}

// printf log
func printfLog(format string, v ...interface{}) {
	if logging == nil {
		return
	}

	logging.Printf(format, v...)
}

// change the name to form filed name
func changeToFieldName(name string) string {
	result := ""
	i := 0
	j := 0
	r := []rune(name)
	for m, v := range r {
		// if the char is the capital
		if v >= 'A' && v < 'a' {
			// if the prior is the lower-case || if the prior is the capital and the latter is the lower-case
			if (m != 0 && r[m-1] >= 'a') || ((m != 0 && r[m-1] >= 'A' && r[m-1] < 'a') && (m != len(r)-1 && r[m+1] >= 'a')) {
				i = j
				j = m
				result += name[i:j] + "_"
			}
		}
	}

	result += name[j:]
	return strings.ToLower(result)
}

// change the params to the query string
func getQueryStr(params interface{}) (result string, err error) {
	if params == nil {
		return
	}
	value := reflect.ValueOf(params)

	switch value.Kind() {
	case reflect.Struct:
		for i := 0; i < value.NumField(); i++ {
			result += "&" + changeToFieldName(value.Type().Field(i).Name) + "=" + fmt.Sprintf("%v", value.Field(i).Interface())
		}
	case reflect.Map:
		for _, key := range value.MapKeys() {
			result += "&" + fmt.Sprintf("%v", key.Interface()) + "=" + fmt.Sprintf("%v", value.MapIndex(key).Interface())
		}
	default:
		err = ErrMustBeStructOrMap
		return
	}

	if result != "" {
		result = result[1:]
	}
	return
}

// sending the request and accept the response, return the response body
func runHandler(req *http.Request) (bodyByte []byte, err error) {

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

// simulate sending file and receive the response, extract the response body
// the first parameter is the method, such as POST,PUT
// the second parameter is the request uri
// the third parameter is the name of the file, containing the directory of the file
// the forth parameter is the field name of the file
// the five parameter is other request params
func TestFileHandler(method, uri, fileName string, fieldName string, param interface{}) (bodyByte []byte, err error) {
	// check whether the router is nil
	if router == nil {
		err = ErrRouterNotSet
		return
	}

	// check whether the method is appropiate, now the method must should be POST or PUT
	if method != POST && method != PUT {
		err = ErrMustPostOrPut
		return
	}

	// create form file
	buf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(buf)

	fileWriter, err := bodyWriter.CreateFormFile(fieldName, fileName)
	if err != nil {
		return
	}

	// open the file
	file, err := os.Open(fileName)
	if err != nil {
		return
	}

	// copy the content of the file to the fileWriter
	length, err := io.Copy(fileWriter, file)
	if err != nil {
		return
	}

	file.Close()
	bodyWriter.Close()

	// make request
	queryStr, err := getQueryStr(param)
	if err != nil {
		return
	}

	if queryStr != "" {
		queryStr = "?" + queryStr
	}

	printfLog("TestFileHandler\tRequest:\t%v:%v%v \tFileName:%v, FieldName:%v\n", method, uri, queryStr, fileName, fieldName)

	req := httptest.NewRequest(method, uri+queryStr, buf)
	req.Header.Add("Content-Type", bodyWriter.FormDataContentType())

	err = req.ParseMultipartForm(length)

	if err != nil {
		return
	}

	// check whether the customed headers are nil
	// if not, then add them to current request headers
	if myHeaders != nil {
		// add the customed headers
		for _, data := range myHeaders {
			if data.IsValid {
				req.Header.Set(data.Key, data.Value)
			}
		}
	}

	// sending request
	bodyByte, err = runHandler(req)

	printfLog("TestFileHandler\tResponse:\t%v:%v,\tResponse:%v\n", method, uri, string(bodyByte))
	return
}

// simulate sending request and receive the response, extract the response body
// the first parameter is the method, such as GET,POST,PUT,DELETE
// the second parameter is the request uri
// the third parameter is the way to send parameter, expected form or json
// the forth parameter is the parameter of the request, it must implement the Param interface
func TestOrdinaryHandler(method string, uri string, way string, param interface{}) (bodyByte []byte, err error) {
	if router == nil {
		err = ErrRouterNotSet
		return
	}

	// make request
	var req *http.Request
	switch way {
	case Json:
		// when the way is Json, change the variable to json bytes, and add to the request body
		jsonBytes := []byte{}
		if param != nil {
			jsonBytes, err = json.Marshal(param)
			if err != nil {
				return
			}
		}
		req = httptest.NewRequest(method, uri, bytes.NewReader(jsonBytes))
		req.Header.Set("Content-Type", "application/json")

		printfLog("TestOrdinaryHandler\tRequest:\t%v:%v,\trequestBody:%v\n", method, uri, string(jsonBytes))
	case Form:
		// when the way is form, then change the variable to query string, and add to the latter of the request uri
		queryStr := ""
		if param != nil {
			queryStr, err = getQueryStr(param)
			if err != nil {
				return
			}
		}

		if queryStr != "" {
			queryStr = "?" + queryStr
		}
		req = httptest.NewRequest(method, uri+queryStr, nil)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		printfLog("TestOrdinaryHandler\tRequest:\t%v:%v%v\n", method, uri, queryStr)
	}

	// check whether the customed headers are nil
	// if not, then add them to current request headers
	if myHeaders != nil {
		// add the customed headers
		for _, data := range myHeaders {
			if data.IsValid {
				req.Header.Set(data.Key, data.Value)
			}
		}
	}

	// launch request
	bodyByte, err = runHandler(req)

	printfLog("TestOrdinaryHandler\tResponse:\t%v:%v\tResponse:%v\n", method, uri, string(bodyByte))
	return
}

// test the ordinary handler and unmarshal the response
func TestHandlerUnMarshalResp(method string, uri string, way string, param interface{}, resp interface{}) error {
	bodyByte, err := TestOrdinaryHandler(method, uri, way, param)
	if err != nil {
		return err
	}

	return json.Unmarshal(bodyByte, resp)
}

// test the file handler and unmarshal the response
func TestFileHandlerUnMarshalResp(method, uri, fileName string, filedName string, param interface{}, resp interface{}) error {
	bodyByte, err := TestFileHandler(method, uri, fileName, filedName, param)
	if err != nil {
		return err
	}

	return json.Unmarshal(bodyByte, resp)
}
