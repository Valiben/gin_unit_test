gin_unit_test
========

A demo about the unit test of gin-gonic/gin

## Installation

Make sure you have a working Go environment (Go 1.2 or higher is required).
See the [install instructions](http://golang.org/doc/install.html).

To install gin_unit_test, simply run:

    go get github.com/Valiben/gin_unit_test

To compile it from source:

    cd $GOPATH/src/github.com/Valiben/gin_unit_test
    go get -u -v
    go build && go test -v

## Example

Here is a simple handler for login. Binding the parameters of the request to the User variable, and judge whether
 the password and username are right.

```go
type User struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password   string `form:"password" json:"password" binding:"required"`
	Age int `form:"age" json:"age" binding:"required"`
}
```
```go
func LoginHandler(c *gin.Context) {
	req := &User{}
	if err := c.Bind(req); err != nil {
		log.Printf("err:%v", err)
		c.JSON(http.StatusOK, gin.H{
			"errno":  "1",
			"errmsg": "parameters not match",
		})
		return
	}

	// judge the password and username
	if req.UserName != "Valiben" || req.Password != "123456" {
		c.JSON(http.StatusOK, gin.H{
			"errno":  "2",
			"errmsg": "password or username is wrong",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"errno":  "0",
		"errmsg": "login success",
	})
}
```

You can write a unit test for this handler like the following.

Firstly, you should set up the router to handle the requests.

```go
router := gin.Default()
router.POST("/login", LoginHandler)
```
Secondly, you should set the router of the utils so that you can use the utils to test the handler.

```go
SetRouter(router)
```
Then you can write the unit test function.

```go
func TestLoginHandler(t *testing.T) {
	resp := OrdinaryResponse{}
	
	err := utils.TestHandlerUnMarshalResp("POST", "/login", "form", user, &resp)
	if err != nil {
		t.Errorf("TestLoginHandler: %v\n", err)
		return
	}
	
	if resp.Errno != "0" {
		t.Errorf("TestLoginHandler: response is not expected\n")
		return
	}
}
````

Then you can run this test to check the handler.

You can find more tests and more specific information about how to use utils in the [test/handlers_test.go](https://github.com/Valiben/gin_unit_test/blob/master/test/handlers_test.go).
