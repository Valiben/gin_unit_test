package test

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"os"
	"testing"
	unitTest "github.com/Valiben/gin_unit_test"
	utils "github.com/Valiben/gin_unit_test/utils"
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
	unitTest.SetRouter(router)

	// set customized request headers
	unitTest.AddHeader(tokenName, myToken)

	newLog := log.New(os.Stdout, "", log.Llongfile|log.Ldate|log.Ltime)
	unitTest.SetLog(newLog)
}

func TestLoginHandler(t *testing.T) {
	// make request params
	param := make(map[string]interface{})
	param["user_name"] = user.UserName
	param["password"] = user.Password
	param["age"] = user.Age

	resp := OrdinaryResponse{}

	err := unitTest.TestHandlerUnMarshalResp("POST", "/login", "form", param, &resp)
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

	err := unitTest.TestHandlerUnMarshalResp("PUT", "/add/user", "form", user, &resp)
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
	err := unitTest.TestHandlerUnMarshalResp(utils.DELETE, "/delete/user", utils.FORM, user, &resp)

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

	resp := OrdinaryResponse{}
	err := unitTest.TestFileHandlerUnMarshalResp(utils.POST, "/upload", (param["file_name"]).(string),
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
	bodyByte, err := unitTest.TestOrdinaryHandler(utils.GET, "/file/test2.txt", utils.FORM, nil)
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
