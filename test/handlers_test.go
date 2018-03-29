package test

import (
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"
	"time"
	utils "zonst/qipai/gin-unittest-demo"
)

func init() {
	// initialize the router
	router := gin.Default()

	// a handler for getting static resources
	router.Static("/file", "./")

	// some handlers for post/put/delete requests
	router.POST("/login", LoginHandler)
	router.PUT("/add/user", AddUserHandler)
	router.DELETE("/delete/user", DeleteUserHandler)
	router.POST("/upload", SaveFileHandler)

	// use a middleware function
	router.Use(Authorize())

	// set the router
	utils.SetRouter(router)

	// set customized request headers
	utils.AddHeader(tokenName, myToken, true)

	newLog := log.New(os.Stdout, "", log.Llongfile|log.Ldate|log.Ltime)
	utils.SetLog(newLog)
}

func TestLoginHandler(t *testing.T) {
	resp := OrdinaryResponse{}

	err := utils.TestHandlerUnMarshalResp(utils.POST, "/login", utils.Form, user, &resp)
	if err != nil {
		t.Errorf("TestLoginHandler: %v\n", err)
		return
	}

	if resp.Errno != "0" {
		t.Errorf("TestLoginHandler: response is not expected\n")
		return
	}
}

func TestAddUserHandler(t *testing.T) {
	resp := OrdinaryResponse{}

	err := utils.TestHandlerUnMarshalResp(utils.PUT, "/add/user", utils.Form, user, &resp)
	if err != nil {
		t.Errorf("TestAddUserHandler: %v\n", err)
		return
	}
	if resp.Errno != "0" {
		t.Errorf("TestAddUserHandler: response is not expected\n")
		return
	}
}

func TestDeleteUserHandler(t *testing.T) {
	resp := OrdinaryResponse{}
	err := utils.TestHandlerUnMarshalResp(utils.DELETE, "/delete/user", utils.Form, user, &resp)

	if err != nil {
		t.Errorf("TestDeleteUserHandler: %v\n", err)
		return
	}

	if resp.Errno != "0" {
		t.Errorf("TestDeleteUserHandler: response is not expected\n")
		return
	}
}

func TestSaveFileHandler(t *testing.T) {

	param := make(map[string]interface{})
	param["file_name"] = "test1.txt"
	param["upload_name"] = "Valiben"

	// get and format the current time string as the content of the next file
	now := time.Now().Format("2006-01-02 15:04:05")

	// create a file
	file, err := os.Create((param["file_name"]).(string))
	if err != nil {
		t.Errorf("TestSaveFileHandler: %v\n", err)
		return
	}

	// use the previous formatted time to initialize a reader
	reader := strings.NewReader(now)

	// copy the content of the reader to the file
	_, err = io.Copy(file, reader)
	if err != nil {
		t.Errorf("TestSaveFileHandler: %v\n", err)
		return
	}

	defer file.Close()

	resp := OrdinaryResponse{}
	err = utils.TestFileHandlerUnMarshalResp(utils.POST, "/upload", (param["file_name"]).(string),
		"file", param, &resp)
	if err != nil {
		t.Errorf("TestSaveFileHandler: %v\n", err)
		return
	}

	if resp.Errno != "0" {
		t.Errorf("TestSaveFileHandler: response is not expected\n")
		return
	}
}

func TestGetFileHandler(t *testing.T) {
	bodyByte, err := utils.TestOrdinaryHandler(utils.GET, "/file/test2.txt", utils.Form, nil)
	if err != nil {
		t.Errorf("TestGetFileHandler: %v\n", err)
		return
	}

	// open the file
	file, err := os.Open("test1.txt")
	if err != nil {
		t.Errorf("TestGetFileHandler: %v\n", err)
		return
	}

	// read all content of the file
	textByte, err := ioutil.ReadAll(file)
	if err != nil {
		t.Errorf("TestGetFileHandler: %v\n", err)
		return
	}

	// judge whether the contents of the two files are equal
	if string(textByte) != string(bodyByte) {
		t.Errorf("TestGetFileHandler: response is not expected\n")
		return
	}
}
