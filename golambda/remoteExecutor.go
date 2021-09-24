package golambda

import (
	lambda "github.com/HarshVaragiya/LambdaFn/liblambda"
	"google.golang.org/grpc"
)

func NewLambdaClient(target string) (lambda.LambdaClient, error) {
	log.Debugf("Attempting to connect to [%s] over gRPC", target)
	conn, err := grpc.Dial(target, grpc.WithInsecure())
	if err != nil {
		log.Warnf("could not connect to server. error = %v", err)
		return nil, err
	}
	client := lambda.NewLambdaClient(conn)
	return client, nil
}
