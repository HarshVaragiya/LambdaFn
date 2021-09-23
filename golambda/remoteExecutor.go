package golambda

import (
	"fmt"
	lambda "github.com/HarshVaragiya/LambdaFn/liblambda"
)

func invokeGrpcRequest(targetIp string, targetPort uint16, targetBootstrapBinary string, event *lambda.Event) (response *lambda.Response, err error) {
	log.Debugf("Remote gRPC invocation to [%s:%v] with bootstrap [%s]", targetIp, targetPort, targetBootstrapBinary)
	return &lambda.Response{}, fmt.Errorf("not implemented")
}
