package main

import (
	"github.com/gin-gonic/gin"
	rvapi "github.com/rivik/go-aux/pkg/rvapi/v2"
)

type apiParams struct {
	apiVer int
}

var problemParseRequest rvapi.ProblemType = rvapi.MustProblemType(rvapi.NewProblemType("problem-types/request-parse-error", "Request parse error"))
var problemGeneralError rvapi.ProblemType = rvapi.MustProblemType(rvapi.NewProblemType("problem-types/general-error", "General error"))

var mainAPI *rvapi.APIServer

func addHandlersV1(r *gin.RouterGroup) {
	r.POST("/example", exampleHandler)
}

func parseParams(c *gin.Context) (*apiParams, error) {
	p := apiParams{}
	p.apiVer = c.MustGet("apiVersion").(int)

	switch p.apiVer {
	default:
	}

	return &p, nil
}

func mainAPIInitAndServeForever(addr string) {
	mainAPI = rvapi.NewAPIServer(addr)

	rvapi.AddHandlersMonitoring(mainAPI.Engine.Group("/"))
	addHandlersV1(mainAPI.Engine.Group("/v1", rvapi.APIVersion(1)))

	mainAPI.ServeForever()
}
