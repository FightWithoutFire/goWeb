package controller

import (
	"github.com/gin-gonic/gin"
	"goWeb/middleware"
	"log"
	"net/http"
)

var (
	GetGinLogo = "/getGinLogo"
)

func init() {
	middleware.GinEngine.GET(GetGinLogo, getGinLogo)
}

func getGinLogo(context *gin.Context) {
	resp, err := http.Get("https://raw.githubusercontent.com/gin-gonic/logo/master/color.png")
	if err != nil {
		log.Println("error get logo")
		context.AbortWithStatusJSON(500, gin.H{
			"msg": "error get logo",
		})

	}

	extraHeader := map[string]string{
		"Content-Disposition": `attachment; filename="gopher.png"`,
	}
	context.DataFromReader(200, resp.ContentLength, resp.Header.Get("Content-Type"), resp.Body, extraHeader)

}
