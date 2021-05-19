package main

import (
	"effective-group-test/cache"
	"github.com/gin-gonic/gin"
	"github.com/pelletier/go-toml"
	"log"
	"net/http"
	"net/http/httputil"
)

const (
	scheme = "https://"
)

func main() {
	config, _ := toml.LoadFile("config.toml")
	backURL := config.Get("main.URL").(string)
	port := config.Get("main.Port").(string)
	maxCountCache := config.Get("Cache.MaxSize").(int64)

	cache := cache.New(uint64(maxCountCache))
	Serve := gin.Default()
	Serve.Use(gin.Recovery())
	any := Serve.Group("/")
	any.GET("/", func(c *gin.Context) {
		body, found := cache.Get(scheme + backURL)
		if !found {
			dump := Dump(scheme + backURL)
			cache.Set(backURL, dump)
			body = &dump
		}
		c.Data(http.StatusOK, "text/html", *body)
	})
	err := Serve.Run(port)
	if err != nil {
		log.Fatal(err)
	}

}
func Dump(Url string) []byte {

	req, err := http.NewRequest("GET", Url, nil)

	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	body, err := httputil.DumpResponse(resp, true)

	if err != nil {
		log.Fatal(err)
	}
	return body
}
