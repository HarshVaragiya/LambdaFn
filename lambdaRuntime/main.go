package main

import (
	"fmt"
	lambda "github.com/HarshVaragiya/LambdaFn/liblambda"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"os"
)

var (
	log     = logrus.New()
	rpcPort = uint16(8888)
	runtime = os.Getenv("LAMBDA_RUNTIME")
	handler = os.Getenv("LAMBDA_HANDLER_FUNCTION")
	bootstrap = os.Getenv("LAMBDA_BOOTSTRAP_BINARY")
	functionName = os.Getenv("LAMBDA_FUNCTION_NAME")
)


func main() {
	log.Printf("Starting gRPC Server on port: %v", rpcPort)
	server := grpc.NewServer()
	rpcServer := lambdaRpcServer{runtime: runtime, handler: handler, bootstrap: bootstrap, functionName: functionName}
	log.Println("registering lambda gRPC server")
	lambda.RegisterLambdaServer(server, rpcServer)
	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", rpcPort))
	if err != nil {
		log.Fatalf("could not start listener. error = %v", err)
	}
	log.Println("attached to listener port. serving gRPC lambda service")
	log.Fatal(server.Serve(listener))
}
