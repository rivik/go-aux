package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	rvapi "github.com/rivik/go-aux/pkg/rvapi/v2"
)

func exampleHandler(c *gin.Context) {
	p, err := parseParams(c)
	if err != nil {
		return
	}
	log.Printf("exampleHandler, api ver: %d", p.apiVer)

	obj := struct {
		ExampleField int `json:"example_field"`
	}{}

	data, err := c.GetRawData()
	if err != nil {
		rvapi.RetErr(c, http.StatusBadRequest, err, problemGeneralError, nil)
		return
	}

	if err := json.Unmarshal(data, &obj); err != nil {
		rvapi.RetErr(c, http.StatusBadRequest, err, problemParseRequest, nil)
		return
	}

	if obj.ExampleField <= 0 {
		rvapi.RetErr(c, http.StatusBadRequest, errors.New("example_field must be > 0"), problemGeneralError, nil)
		return
	}

	log.Printf("Request: %+v\n", obj)
	rvapi.Ret(c, http.StatusOK, &obj)
}
