package test

import (
	"net/http"
	"github.com/gin-gonic/gin/binding"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"io"
	"fmt"
)

type User struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password   string `form:"password" json:"password" binding:"required"`
	Age int `form:"age" json:"age" binding:"required"`
}

// Login
// 登录请求处理函数
func Login(c *gin.Context) {
	req := &User{}
	if err := c.ShouldBindWith(req, binding.JSON); err != nil {
		log.Printf("err:%v",err)
		c.JSON(http.StatusOK, gin.H{
			"errno":"1",
			"errmsg":"parameters not match",
		})
		return
	}

	// judge the password and username
	// 判断密码与用户名是否匹配
	if req.Username != "Valiben" || req.Password != "123456" {
		c.JSON(http.StatusOK, gin.H{
			"errno":"2",
			"errmsg":"password or username is wrong",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"errno":"0",
		"errmsg":"login success",
	})
}

// DeleteUser
// 删除用户信息请求处理函数
func DeleteUser(c *gin.Context) {
	req := &User{}
	if err := c.ShouldBindWith(req, binding.JSON); err != nil {
		log.Printf("err:%v",err)
		c.JSON(http.StatusOK, gin.H{
			"errno":"1",
			"errmsg":"parameters not match",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"errno":"0",
		"errmsg":fmt.Sprintf("delete user:%+v", req),
	})
}

// AddUser
// 添加用户信息请求处理函数
func AddUser(c *gin.Context) {
	req := &User{}
	if err := c.ShouldBindWith(req, binding.JSON); err != nil {
		log.Printf("err:%v",err)
		c.JSON(http.StatusOK, gin.H{
			"errno":"1",
			"errmsg":"parameters not match",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"errno":"0",
		"errmsg":fmt.Sprintf("add user:%+v", req),
	})
}

// SaveFile
// 保存文件
func SaveFile(c *gin.Context) {
	// get the file of the request
	// 获取请求中的文件
	file,_,_ := c.Request.FormFile("file")
	if file == nil {
		c.JSON(http.StatusOK, gin.H{
			"errno":"2",
			"errmsg":"file is nil",
		})
		return
	}
	out,err := os.Create("test2.txt")

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"errno":"2",
			"errmsg":err.Error(),
		})
		return
	}

	// copy the content of the file to the out
	// 将获取到的文件中的内容拷贝到out输出流对应的文件
	_,err = io.Copy(out, file)
	defer file.Close()
	defer out.Close()

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"errno":"2",
			"errmsg":err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"errno":"0",
		"errmsg":"save file success",
	})
}