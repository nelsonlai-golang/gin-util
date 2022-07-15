package rest

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/nelsonlai-golang/go-util/random"
)

// response is a REST response object
type response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

type errorMsg struct {
	Ref     string `json:"ref"`
	Message string `json:"msg"`
}

// OK returns a successful response
func OK(c *gin.Context, data interface{}) {
	c.JSON(200, response{Success: true, Data: data})
}

// FAIL returns a failed response
func FAIL(c *gin.Context, err error) {
	ref := random.RandomString(10, random.StringConfig{
		Uppercase: true,
		Lowercase: true,
	})
	log.Default().Println("ref: ", ref, " | err: ", err.Error())
	c.JSON(200, response{Success: false, Data: errorMsg{Ref: ref, Message: err.Error()}})
}
