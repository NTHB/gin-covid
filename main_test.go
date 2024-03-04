package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/stretchr/testify/assert"
)

func SetUpRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to covid summary project!",
		})
	})
	return router
}

func TestHomepage(t *testing.T) {
	mockResponse := `{"message":"Welcome to covid summary project!"}`
	r := SetUpRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, http.StatusOK, w.Code)
	responseData, _ := ioutil.ReadAll(w.Body)
	assert.Equal(t, mockResponse, string(responseData))
}
