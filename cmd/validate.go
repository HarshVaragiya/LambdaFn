package cmd

import (
	"fmt"
	"github.com/HarshVaragiya/LambdaFn/golambda"
)

func ValidateCreateLambda(function *golambda.Function) error {
	if function.Name == "" || function.Handler == "" || function.CodeUri == "" || function.Runtime == "" {
		return fmt.Errorf("cannot create function with empty fields in [name, handler, code-uri, runtime]")
	}
	return nil
}
