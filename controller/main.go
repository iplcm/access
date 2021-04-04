package controller

import (
	"bytes"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/iplcm/access/tools"
)

var (
	ApiHost = "https://www.aloy.asia"
)

func Info(c *gin.Context) {
	session := sessions.Default(c)
	url := ApiHost + "/capi/info"
	user, passwd := session.Get("user"), session.Get("passwd")
	reqest, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(fmt.Sprintf(`{"user":"%s","passwd":"%s"}`, user, passwd))))
	type Res struct {
		Ret  int                    `json:"ret"`
		Data map[string]interface{} `json:"data"`
	}
	var res Res
	if curlerr := tools.CURL(reqest, &res); curlerr == nil && err == nil && res.Ret == 1 {
		// 无报错且有结果返回
		uuri := make(map[string]string)
		for k, v := range res.Data["uuri"].(map[string]interface{}) {
			uuri[k] = v.(string)
		}
		including := []string{"Shadowsocks", "Clash", "Surge3", "Surge4", "Surfboard", "Kitsunebi", "QuantumultX", "Quantumult", "Shadowrocket", "v2rayN/G", "SurgeNode", "ClashNode"}
		sort.Strings(including)
		c.HTML(http.StatusOK, "info", gin.H{
			"text":      strings.ReplaceAll(res.Data["text"].(string), "\n", "<br>"),
			"including": including,
			"uuri":      uuri,
			"host":      res.Data["host"].(string),
			"name":      res.Data["name"].(string),
			"email":     res.Data["email"].(string),
		})
	} else if curlerr == nil && res.Ret == 0 {
		// 异常登出账户
		session.Clear()
		_ = session.Save()
	} else {
		c.String(http.StatusOK, "服务器异常")
	}
}

func Sync(c *gin.Context) {
	url := ApiHost + "/capi/sync"
	session := sessions.Default(c)
	reqest, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(fmt.Sprintf(`{"user":"%s","passwd":"%s"}`, session.Get("user"), session.Get("passwd")))))
	var ret map[string]string
	if curlerr := tools.CURL(reqest, &ret); curlerr == nil && err == nil {
		if ret["ret"] == "1" {
			tools.Redirect(http.StatusFound, ret["host"]+fmt.Sprintf("/capi/sync?token=%s", ret["token"]), c)
		} else {
			Logout(c)
		}
	} else {
		c.String(http.StatusOK, "ERROR")
	}
}

func Login(c *gin.Context) {
	session := sessions.Default(c)
	if session.Get("is_login") == true {
		tools.Redirect(http.StatusFound, "/get/info", c)
		return
	}
	returnLoginHtml(c, nil)
}

func Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	_ = session.Save()
	tools.Redirect(http.StatusFound, "/login", c)
}

func returnLoginHtml(c *gin.Context, tips interface{}) {
	var hidden string
	if tips == nil {
		hidden = "hidden"
	}
	c.HTML(http.StatusOK, "login", gin.H{
		"tips":   tips,
		"hidden": hidden,
	})
}

func DoLogin(c *gin.Context) {
	url := ApiHost + "/capi/login"
	user, passwd := c.PostForm("user"), c.PostForm("passwd")
	reqest, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(fmt.Sprintf(`{"user":"%s","passwd":"%s"}`, user, passwd))))
	var ret map[string]int
	if curlerr := tools.CURL(reqest, &ret); curlerr == nil && err == nil {
		if ret["ret"] == 1 {
			session := sessions.Default(c)
			session.Set("is_login", true)
			session.Set("user", user)
			session.Set("passwd", passwd)
			_ = session.Save()
			tools.Redirect(http.StatusFound, "/get/info", c)
		} else {
			returnLoginHtml(c, "账户或密码错误")
		}
	} else {
		returnLoginHtml(c, "服务器异常")
	}
}
