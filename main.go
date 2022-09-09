package main

import (
	"flag"
	"net/http"
	"net/http/httputil"
	"net/url"

	_ "embed"

	"github.com/gin-gonic/gin"
)

var fPort = flag.String("port", "8080", "which port will be used")
var port string

var fApiHost = flag.String("api", "http://localhost:3000", "base url whatsapp api")
var apiHost string

//go:embed src/index.html
var indexHTML string

//go:embed src/app.css
var appCSS string

//go:embed src/app.js
var appJS string

func init() {
	flag.Parse()
	port = *fPort
	apiHost = *fApiHost
}

func htmlRoute(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(indexHTML))
	})
	r.GET("/app.css", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/css; charset=utf-8", []byte(appCSS))
	})
	r.GET("/app.js", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/javascript; charset=utf-8", []byte(appJS))
	})
}

func apiRoute(r *gin.Engine) {
	r.Any("/api/*proxyPath", func(c *gin.Context) {
		remote, err := url.Parse(apiHost)
		if err != nil {
			panic(err)
		}
		proxy := httputil.NewSingleHostReverseProxy(remote)
		proxy.Director = func(req *http.Request) {
			req.Header = c.Request.Header
			req.Host = remote.Host
			req.URL.Scheme = remote.Scheme
			req.URL.Host = remote.Host
			req.URL.Path = "/api" + c.Param("proxyPath")
		}
		proxy.ServeHTTP(c.Writer, c.Request)
	})
}

func main() {
	r := gin.Default()
	r.NoRoute(func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(indexHTML))
	})

	htmlRoute(r)
	apiRoute(r)

	r.Run(":" + port)
}
