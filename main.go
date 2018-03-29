package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

var gmina *GMina

func main() {
	var (
		apiKey     string
		cacheDir   string
		listenAddr string
		debugMode  bool

		err error
	)
	flag.StringVar(&apiKey, "apikey", "", "Google API key")
	flag.StringVar(&cacheDir, "cache", "./cache", "Cache directory")
	flag.StringVar(&listenAddr, "addr", "localhost:8585", "Listen address")
	flag.BoolVar(&debugMode, "debug", false, "GIN debug mode")

	flag.Parse()

	if apiKey == "" {
		fmt.Println("API key is empty.")
		flag.Usage()
		return
	}

	gmina, err = New(apiKey, cacheDir)
	if err != nil {
		fmt.Println(err)
		return
	}

	if !debugMode {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	router.GET("/tts", ttsHandler)

	router.Run(listenAddr)
}

func ttsHandler(ctx *gin.Context) {
	var text = ctx.Query("text")
	content, err := gmina.TTS(text)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "bad request", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "ok", "message": "everything is good !", "content": content})
}
