package main

import (
	"fmt"
	"github.com/HarshVaragiya/LambdaFn/golambda"
	"github.com/HarshVaragiya/LambdaFn/liblambda"
	"sync"
	"time"
)

type LambdaManager struct {
	lock             *sync.RWMutex
	awsAccountNumber int64 // used to generate function ARN
	lambdaFunctions  map[string]golambda.Function
}

func NewLambdaManager(awsAccountNumber int64) *LambdaManager {
	lock := &sync.RWMutex{}
	return &LambdaManager{awsAccountNumber: awsAccountNumber, lock: lock}
}

func (manager *LambdaManager) CreateLambdaFunction(function *golambda.Function) (*liblambda.LambdaRestApiResponse, error) {
	var err error
	manager.lock.Lock() // ensure we don't do anything else concurrently while creating new function
	defer manager.lock.Unlock()
	if _, exists := manager.lambdaFunctions[function.Name]; exists {
		errMsg := fmt.Sprintf("function with name [%s] already exists. cannot create another function with same name.", function.Name)
		log.Warn(errMsg)
		return liblambda.NewErrorRestResponse(errMsg, function.Name, "function exists"), fmt.Errorf("function exists")
	}
	if function.Timeout, err = time.ParseDuration(function.TimeoutSeconds); err != nil {
		errMsg := fmt.Sprintf("error parsing time for new lambda function [%s]. time string : [%s]. error = %v", function.Name, function.TimeoutSeconds, err)
		log.Warn(err)
		return liblambda.NewErrorRestResponse(errMsg, function.Name, err.Error()), err
	}
	function.Executor = golambda.NewContainerExecutor(function)
	function.CalculateArn(manager.awsAccountNumber)
	msg := fmt.Sprintf("created function [%s] with arn [%s]", function.Name, function.Arn)
	log.Info(msg)
	return liblambda.NewOkRestResponse(msg, function.Name, function.Arn, msg), nil
}

func (manager *LambdaManager) RemoveLambdaFunction(function *golambda.Function) (*liblambda.LambdaRestApiResponse, error) {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	mappedFunction, exists := manager.lambdaFunctions[function.Name]
	if !exists {
		errMsg := fmt.Sprintf("attempting to delete a non-existent lambda function")
		debugMsg := fmt.Sprintf("cannot delete lambda function with name [%s] as it does not exist", function.Name)
		log.Warn(debugMsg)
		return liblambda.NewErrorRestResponse(errMsg, function.Name, debugMsg), fmt.Errorf("function does not exist")
	}
	msg := fmt.Sprintf("removing lambda function [%s]", mappedFunction.Arn)
	debugMsg := fmt.Sprintf("removing lambda function with name [ %s ], and arn [ %s ]", mappedFunction.Name, mappedFunction.Arn)
	delete(manager.lambdaFunctions, function.Name)
	log.Println(debugMsg)
	return liblambda.NewOkRestResponse(msg, function.Name, mappedFunction.Arn, debugMsg), nil
}

func (manager *LambdaManager) InvokeLambdaFunction(functionName string, event *liblambda.Event) (*liblambda.Response, error) {
	manager.lock.RLock()
	mappedFunction, exists := manager.lambdaFunctions[functionName]
	manager.lock.RUnlock()
	if !exists {
		errMsg := fmt.Sprintf("attempting to invoke a non-existent lambda function")
		debugMsg := fmt.Sprintf("cannot invoke lambda function with name [%s] as it does not exist", functionName)
		log.Warn(debugMsg)
		return &liblambda.Response{Message: errMsg, StatusCode: statusCodeMap["does-not-exist"], Stderr: debugMsg}, fmt.Errorf("function does not exist")
	}
	log.Infof("invoking function [%s] with eventId %s", functionName, event.EventId)
	return mappedFunction.Invoke(event)
}

func (manager *LambdaManager) ModifyLambdaFunction(function *golambda.Function) (*liblambda.LambdaRestApiResponse, error) {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	mappedFunction, exists := manager.lambdaFunctions[function.Name]
	if !exists {
		errMsg := fmt.Sprintf("attempting to modify a non-existent lambda function")
		debugMsg := fmt.Sprintf("cannot modify lambda function with name [%s] as it does not exist", function.Name)
		log.Warn(debugMsg)
		return liblambda.NewErrorRestResponse(errMsg, function.Name, debugMsg), fmt.Errorf("function does not exist")
	}
	errMsg := fmt.Sprintf("error modifying lambda function with arn [%s]. error = NOT IMPLEMENTED", mappedFunction.Arn)
	return liblambda.NewErrorRestResponse(errMsg, function.Name, errMsg), fmt.Errorf("not implemented")
}

func (manager *LambdaManager) GetManagedLambdaFunctions() map[string]string {
	manager.lock.RLock()
	defer manager.lock.RUnlock()
	managedFunctions := make(map[string]string, len(manager.lambdaFunctions))
	for fname, function := range manager.lambdaFunctions {
		managedFunctions[fname] = function.Arn
	}
	return managedFunctions
}
