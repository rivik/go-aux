package rvapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rivik/go-aux/pkg/appver"
)

const (
	HeaderContentType = "Content-Type"
)

type StatusResponse struct {
	Status     int    `json:"status,omitempty"`
	StatusText string `json:"status_text,omitempty"`
}

type PingResponse struct {
	IsAlive bool              `json:"is_alive"`
	AppVer  appver.AppVersion `json:"app_version"`
}

func Ret(c *gin.Context, code int, data interface{}) {
	c.JSON(code, data)
}

func RetErr(c *gin.Context, code int, err error, problemType ProblemType, extensions interface{}) {
	c.Header(HeaderContentType, MediaTypeProblemJSON)

	pd := NewProblemDetails(
		problemType,
		code,
		err.Error(),
		c.FullPath(),
		extensions,
	)
	Ret(c, code, pd)
}

func RetStatusText(c *gin.Context, code int) {
	Ret(c, code, StatusResponse{Status: code, StatusText: http.StatusText(code)})
}
