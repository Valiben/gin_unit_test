package utils

import (
	"net/http"
	"bytes"
	"encoding/json"
	"errors"
	"mime/multipart"
	"strings"
	"io/ioutil"
)

const (
	GET = "GET"
	POST = "POST"
	PUT = "PUT"
	DELETE = "DELETE"

	JSON = "json"
	FORM = "form"
)

var(
	ErrMethodNotSupported     = errors.New("method is not supported")
	ErrMIMENotSupported = errors.New("mime is not supported")
)

// make request which contains uploading file
func MakeFileRequest(method, api, fileName,fieldName string, param interface{}) (request *http.Request, err error) {
	method = strings.ToUpper(method)
	if method != POST && method != PUT {
		err = ErrMethodNotSupported
		return
	}

	// create form file
	buf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(buf)
	fileWriter, err := bodyWriter.CreateFormFile(fieldName, fileName)
	if err != nil {
		return
	}

	// read the file
	fileBytes,err := ioutil.ReadFile(fileName)
	if err != nil {
		return
	}

	// read the file to the fileWriter
	length, err := fileWriter.Write(fileBytes)
	if err != nil {
		return
	}

	bodyWriter.Close()

	// make request
	queryStr := MakeQueryStrFrom(param)
	if queryStr != "" {
		api += "?"+queryStr
	}
	request,err = http.NewRequest(string(method), api, buf)
	if err != nil {
		return
	}

	request.Header.Set("Content-Type", bodyWriter.FormDataContentType())
	err = request.ParseMultipartForm(int64(length))
	return
}

// make request
func MakeRequest(method, mime, api string, param interface{}) (request *http.Request, err error) {
	method = strings.ToUpper(method)
	mime = strings.ToLower(mime)

	switch mime {
	case JSON:
		var (
			contentBuffer *bytes.Buffer
			jsonBytes []byte
		)
		jsonBytes, err = json.Marshal(param)
		if err != nil {
			return
		}
		contentBuffer = bytes.NewBuffer(jsonBytes)
		request,err = http.NewRequest(string(method), api, contentBuffer)
		if err != nil {
			return
		}
		request.Header.Set("Content-Type", "application/json;charset=utf-8")
	case FORM:
		queryStr := MakeQueryStrFrom(param)
		buffer := bytes.NewReader([]byte(queryStr))
		if ( method == DELETE || method == GET ) && queryStr != "" {
			api += "?"+queryStr
		}
		request,err = http.NewRequest(string(method), api, buffer)
		if err != nil {
			return
		}
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")
	default:
		err = ErrMIMENotSupported
		return
	}
	return
}
