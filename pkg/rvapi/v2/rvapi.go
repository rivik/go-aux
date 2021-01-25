package rvapi

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rivik/go-aux/pkg/appver"
)

func ApiVersion(ver int) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("apiVersion", ver)
		c.Next()
	}
}

func PingMe(c *gin.Context) {
	Ret(c, http.StatusOK, PingResponse{IsAlive: true, AppVer: appver.Version})
}

func AddHandlersMonitoring(r *gin.RouterGroup) {
	r.GET("/ping", PingMe)
}

type APIServer struct {
	Engine *gin.Engine
	Server *http.Server
}

func NewAPIServer(addr string) *APIServer {
	api := APIServer{Engine: gin.New()}

	api.Engine.HandleMethodNotAllowed = true

	api.Server = &http.Server{
		Addr:    addr,
		Handler: api.Engine,
	}

	return &api
}

func (api *APIServer) ServeForever() {
	log.Printf("Starting api server at %s ...\n", api.Server.Addr)
	log.Fatal(api.Server.ListenAndServe())
}
