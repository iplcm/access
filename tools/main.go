package tools

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func NewUUID() uuid.UUID {
	uuid, _ := uuid.NewRandom()
	return uuid
}
func Redirect(code int, url string, c *gin.Context) {
	c.Redirect(code, url)
	c.Abort()
}

func CURL(r *http.Request, interfaces interface{}) error {
	r.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(r)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if body, err := ioutil.ReadAll(response.Body); err == nil {
		bytes := []byte(body)
		return json.Unmarshal(bytes, interfaces)
	} else {
		return err
	}
}
func Unescaped(x string) interface{} {
	return template.HTML(x)
}
