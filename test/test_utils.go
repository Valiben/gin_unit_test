package test

import (
	"net/http"
	"fmt"
	"encoding/json"
	"github.com/gin-gonic/gin"
)

var (
	// a default user variable for the next requests parameter
	// 作为接下来多个请求的参数
	user = &User{
		Username:"Valiben",
		Password:"123456",
		Age:22,
	}

	// a default token for authorization
	// 默认的token验证
	myToken = "ssoiuoiu"
)

// receive the ordinary response
// 接收普通响应的结构体
type OrdinaryResponse struct {
	Errno string `json:"errno"`
	Errmsg string `json:"errmsg"`
}

// implements the Param.QueryStr function, make a query string for User variable
// 实现了Param接口的QueryStr方法，用于将User结构体的各成员变量构造成一个query string
func (u *User) QueryStr() string {
	return fmt.Sprintf("username=%v&password=%v&age=%v", u.Username, u.Password, u.Age)
}

// implements the Param.JsonBytes function, make a json []byte for User variable
// 实现了Param接口的JsonBytes方法，用于将User结构体的各成员变量构造成一个json类型的[]byte
func (u *User) JsonBytes() []byte {
	jsonByte,_ := json.Marshal(u)
	return jsonByte
}

// a middleware function for easy authorization
// 用于鉴权的中间件
func Authorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("x-xq5-jwt")
		if token != myToken {
			// unauthorized, just return the unauthorized response
			// 鉴权失败，拦截请求，直接给出响应
			c.JSON(http.StatusUnauthorized, gin.H{"errno": "-1", "errmsg": "unauthorized"})
			return
		}else {
			// authorized, come into the next function
			// 鉴权成功，进入接下来的请求处理函数
			c.Next()
		}
	}
}

