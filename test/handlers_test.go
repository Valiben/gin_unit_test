package test

import (
	"testing"
	"fmt"
	"github.com/gin-gonic/gin"
	utils "github.com/Valiben/gin_unit_test"
	"encoding/json"
	"time"
	"os"
	"strings"
	"io"
	"io/ioutil"
)

func init() {
	// initialize the router
	// 初始化路由
	router := gin.Default()

	// use a middleware function
	// 使用鉴权函数作为中间件
	router.Use(Authorize())

	// a handler for getting static resources
	// 访问静态资源的路由
	router.Static("/file", "./")

	// some handlers for post/put/delete requests
	// 一些用于处理post,put,delete的handler
	router.POST("/login", Login)
	router.PUT("/add/user", AddUser)
	router.DELETE("/delete/user", DeleteUser)
	router.POST("/upload", SaveFile)

	// set the router
	// 设置util中的router，以便utils中的工具函数调用
	utils.SetRouter(router)

	// customize some request headers, such as x-xq5-jwt for authorization
	// 自定义请求头，用于传递token，以便鉴权
	headers := make(map[string]string)
	headers["x-xq5-jwt"] = myToken

	// set customized request headers
	// 设置自定义请求头，以便utils中的工具函数调用
	utils.SetCustomedHeaders(headers)
}

// TestLogin
func TestLogin(t *testing.T) {
	// make the request params
	// 构造请求需要传递的参数
	egParam := &utils.EgParam{
		make(map[string]interface{}),
	}
	egParam.Set("username", user.Username)
	egParam.Set("password", user.Password)
	egParam.Set("age", user.Age)

	// simulate sending ordinary request
	// 模拟发送请求
	bodyByte,err := utils.SimulateOrdinaryHandler(utils.POST, "/login", utils.Json, egParam)
	if err != nil {
		t.Errorf("TestLogin: %v\n", err)
		return
	}
	fmt.Printf("TestLogin: response: %v\n", string(bodyByte))

	// make a variable to receive the response body bytes
	// 用于接收响应body中的json bytes
	resp := OrdinaryResponse{}
	if err := json.Unmarshal(bodyByte, &resp); err != nil {
		t.Errorf("TestLogin: %v\n", err)
		return
	}

	if resp.Errno != "0" {
		t.Errorf("TestLogin: response is not expected\n")
		return
	}
}

func TestAddUser(t *testing.T) {
	bodyByte,err := utils.SimulateOrdinaryHandler(utils.PUT, "/add/user", utils.Json, user)
	if err != nil {
		t.Errorf("TestAddUser: %v\n", err)
		return
	}
	fmt.Printf("TestAddUser: response: %v\n", string(bodyByte))

	resp := OrdinaryResponse{}
	if err := json.Unmarshal(bodyByte, &resp); err != nil {
		t.Errorf("TestAddUser: %v\n", err)
		return
	}

	if resp.Errno != "0" {
		t.Errorf("TestAddUser: response is not expected\n")
		return
	}
}

func TestDeleteUser(t *testing.T) {
	bodyByte,err := utils.SimulateOrdinaryHandler(utils.DELETE, "/delete/user", utils.Json, user)
	if err != nil {
		t.Errorf("TestDeleteUser: %v\n", err)
		return
	}
	fmt.Printf("TestDeleteUser: response: %v\n", string(bodyByte))

	resp := OrdinaryResponse{}
	if err := json.Unmarshal(bodyByte, &resp); err != nil {
		t.Errorf("TestDeleteUser: %v\n", err)
		return
	}

	if resp.Errno != "0" {
		t.Errorf("TestDeleteUser: response is not expected\n")
		return
	}
}

func TestSaveFile(t *testing.T) {
	// get and format the current time string as the content of the next file
	// 获取并初始化现在的时间，并将其作为下面文件的内容
	now := time.Now().Format("2006-01-02 15:04:05")

	// create a file
	// 创建一个文件
	file,err := os.Create("test1.txt")
	if err != nil {
		t.Errorf("TestSaveFile: %v\n", err)
		return
	}

	// use the previous formatted time to initialize a reader
	// 使用先前格式化好的时间string初始化一个reader
	reader := strings.NewReader(now)

	// copy the content of the reader to the file
	// 将reader中的内容拷贝到file中
	_,err = io.Copy(file, reader)
	if err != nil {
		t.Errorf("TestSaveFile: %v\n", err)
		return
	}

	defer file.Close()

	bodyByte,err := utils.SimulateFileHandler(utils.POST, "/upload", "test1.txt", "file")
	if err != nil {
		t.Errorf("TestSaveFile: %v\n", err)
		return
	}
	fmt.Printf("TestSaveFile: response: %v\n", string(bodyByte))

	resp := OrdinaryResponse{}
	if err := json.Unmarshal(bodyByte, &resp); err != nil {
		t.Errorf("TestSaveFile: %v\n", err)
		return
	}

	if resp.Errno != "0" {
		t.Errorf("TestSaveFile: response is not expected\n")
		return
	}
}

func TestGetFile(t *testing.T) {
	bodyByte,err := utils.SimulateOrdinaryHandler(utils.GET, "/file/test2.txt", utils.Form, nil)
	if err != nil {
		t.Errorf("TestGetFile: %v\n", err)
		return
	}

	// open the file
	// 打开文件
	file,err := os.Open("test1.txt")
	if err != nil {
		t.Errorf("TestGetFile: %v\n", err)
		return
	}

	// read all content of the file
	// 读出文件的所有内容
	textByte,err := ioutil.ReadAll(file)
	if err != nil {
		t.Errorf("TestGetFile: %v\n", err)
		return
	}

	fmt.Printf("TestGetFile: bodyByte: %v, textByte: %v\n", string(bodyByte), string(textByte))

	// judge whether the contents of the two files are equal
	// 比较两个文件的内容是否相等
	if string(textByte) != string(bodyByte) {
		t.Errorf("TestGetFile: response is not expected\n")
		return
	}
}