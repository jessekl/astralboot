// Serve web interface
//go:generate go-bindata -debug asset/...
package main

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"html/template"
	"io"
	"strconv"
	"strings"
)

// WebInterface : provides a web interface for astralboot functions and monitoring
func (wh *WebHandler) WebInterface() {
	// Bind the Index
	wh.router.GET("/", wh.Index)
	wh.router.GET("/static/*path", wh.Static)
	// Load the templates
	// get the asset dir
	pages, err := AssetDir("asset/pages")
	if err != nil {
		logger.Error("Loading pages %s", err)
		return
	}
	templ := template.New("")
	for _, j := range pages {
		logger.Critical("%s", j)
		data, _ := Asset(j)
		fmt.Println(data)
		tmpl.New(j).Parse(string(data))
	}
	wh.uiTemplates = templ

}

func (wh *WebHandler) Index(c *gin.Context) {
	logger.Debug("Index HIT")
	wh.uiTemplates.ExecuteTemplate(c.Writer, "index.html", nil)
}

func (wh *WebHandler) Static(c *gin.Context) {
	path := c.Params.ByName("path")
	logger.Debug(path)
	data, err := Asset("asset" + path)
	if err != nil {
		logger.Error("Asset Error ", err)
		c.AbortWithStatus(404)
	}
	if strings.HasSuffix(path, ".css") {
		c.Writer.Header().Set("Content-Type", "text/css")
	}
	if strings.HasSuffix(path, ".js") {
		c.Writer.Header().Set("Content-Type", "text/javascript")
	}
	size := int64(len(data))
	c.Writer.Header().Set("Content-Length", strconv.FormatInt(size, 10))
	io.Copy(c.Writer, bytes.NewReader(data))
}
