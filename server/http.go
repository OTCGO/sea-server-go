package server

import (
	"github.com/OTCGO/sea-server-go/config"
	"github.com/OTCGO/sea-server-go/service"
	"github.com/gin-gonic/gin"
	"github.com/hzxiao/goutil"
	"github.com/hzxiao/goutil/log"
	"net/http"
)

var g *gin.Engine
var neoService *service.NeoService

func Init()  {
	g = gin.New()
	neoService = service.NewNeoService()
}

func Run() {
	g.Use(gin.Recovery())
	g.Use(NoCache)
	g.Use(Options)
	g.Use(Cors())
	// 404 Handler.
	g.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "The incorrect API route.")
	})

	registerHandler(g)

	log.Info("[Server] listen on %v", config.Conf.Port)
	log.Info(http.ListenAndServe(config.Conf.Port, g).Error())
}

func registerHandler(g *gin.Engine) {
	net := g.Group("/:net")

	adm := net.Group("/adm")
	{
		adm.GET("/:height/:task/mockSync", mockSync)
		adm.GET("/stats", stats)
	}

	net.GET("/height", neoHeight)
	net.GET("/block/:height", neoBlock)
}

func WriteJSON(c *gin.Context, data interface{}, err error) {
	if err != nil {
		log.Error("http path: %v, error: %v", c.Request.URL.Path, err)
	}

	c.JSON(http.StatusOK, formatResult(data, err))
}

func formatResult(data interface{}, err error) interface{} {
	if err != nil {
		if data == nil {
			return goutil.Map{"result": false, "error": err.Error()}
		}
		if m, ok := data.(goutil.Map); ok {
			m.Set("result", false)
			m.Set("error", err.Error())
			return m
		}
	}

	if data == nil {
		return data
	}
	m, ok := data.(goutil.Map)
	if ok {
		m.Set("result", true)
		return m
	}

	return data
}
