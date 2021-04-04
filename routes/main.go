package routes

import (
	"html/template"
	"net/http"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/iplcm/access/controller"
	"github.com/iplcm/access/tools"
)

type Routes struct {
	routes *gin.Engine
}

func Init() Routes {
	r := Routes{routes: gin.Default()}
	r.routes.SetFuncMap(template.FuncMap{
		"Unescaped": tools.Unescaped,
	})
	r.routes.LoadHTMLGlob("template/**/*.tmpl")
	r.routes.Static("/static", "template/login-form-20")
	return r
}

func (r Routes) Run() {
	_ = r.routes.RunTLS("0.0.0.0:2096", "ssl/test.pub", "ssl/test.pem")
}

func (r Routes) SetSession() Routes {
	// go run cli/main.go uuid
	r.routes.Use(sessions.Sessions("crm-session", cookie.NewStore([]byte("96d0965e-e451-4582-9e8b-9aa19d8e446e"))))
	return r
}
func (r Routes) AddRoutes() Routes {
	site := r.routes.Group("", func(c *gin.Context) {
		c.Header("x-powered-by", "@CR2477")
		if c.Request.Host != os.Getenv("CRM_HOST") {
			c.AbortWithStatus(444)
		}
		c.Next()
	})
	site.GET("/", func(c *gin.Context) { c.Redirect(http.StatusFound, "/login") })
	site.GET("/login", controller.Login)
	site.GET("/logout", controller.Logout)
	site.POST("/login", controller.DoLogin)
	get := site.Group("/get", func(c *gin.Context) {
		session := sessions.Default(c)
		if session.Get("is_login") != true {
			tools.Redirect(http.StatusFound, "/login", c)
			return
		}
		c.Next()
	})
	get.GET("/info", controller.Info)
	get.GET("/sync", controller.Sync)
	return r
}
