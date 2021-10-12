package cmd

import (
	"github.com/HarshVaragiya/LambdaFn/golambda"
	"github.com/HarshVaragiya/LambdaFn/liblambda"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type LambdaRestApi struct {
	lambdaManager *LambdaManager
	router        *gin.Engine
}

func (api *LambdaRestApi) Init() {
	api.router.PUT("/lambda", api.Create)
	api.router.GET("/lambda/:functionName/invoke", api.Invoke)
	api.router.DELETE("/lambda", api.Delete)
}

func (api *LambdaRestApi) Create(ctx *gin.Context) {
	var function golambda.Function
	if err := ctx.BindJSON(&function); err == nil {
		log.Infof("creating new lambda function [%s]", function.Name)
		resp, err := api.lambdaManager.CreateLambdaFunction(&function)
		respCode := http.StatusCreated
		if err != nil {
			log.Warnf("error while creating new lambda function [%s]. error = %v", function.Name, err)
			respCode = http.StatusBadRequest
		}
		ctx.IndentedJSON(respCode, resp)
	} else {
		log.Warnf("error parsing create request. error = %v", err)
		ctx.AbortWithStatus(http.StatusBadRequest)
	}
}

func (api *LambdaRestApi) Invoke(ctx *gin.Context) {
	functionName := ctx.Param("functionName")
	var event liblambda.Event
	log.Debugf("received lambda invocation request for lambda [%s]", functionName)
	if err := ctx.BindJSON(&event); err == nil {
		event.EventId = uuid.NewString()
		log.Infof("invoking lambda [%s] with eventId [%s]", functionName, event.EventId)
		resp, err := api.lambdaManager.InvokeLambdaFunction(functionName, &event)
		respCode := http.StatusOK
		if err != nil {
			log.Warnf("error invoking lambda function [%s] . error = %v", functionName, err)
			respCode = http.StatusFailedDependency
		}
		ctx.IndentedJSON(respCode, resp)
	} else {
		log.Warnf("error parsing invocation request for lambda [%v]. error = %v", functionName, err)
		ctx.AbortWithStatus(http.StatusBadRequest)
	}
}

func (api *LambdaRestApi) Delete(ctx *gin.Context) {
	var function golambda.Function
	if err := ctx.BindJSON(&function); err == nil {
		log.Infof("removing lambda function [%s] ", function.Name)
		resp, err := api.lambdaManager.RemoveLambdaFunction(&function)
		respCode := http.StatusOK
		if err != nil {
			log.Warnf("error while removing lambda function [%s]. error = %v", function.Name, err)
			respCode = http.StatusNotFound
		}
		ctx.IndentedJSON(respCode, resp)
	} else {
		log.Warnf("error parsing delete function request. error = %v", err)
		ctx.AbortWithStatus(http.StatusBadRequest)
	}
}
