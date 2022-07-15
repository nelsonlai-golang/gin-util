package rest

import (
	"log"
	"math/rand"

	"github.com/gin-gonic/gin"
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
	ref := _randStringRunes(10)
	log.Default().Println("ref: ", ref, " | err: ", err.Error())
	c.JSON(200, response{Success: false, Data: errorMsg{Ref: ref, Message: err.Error()}})
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func _randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
