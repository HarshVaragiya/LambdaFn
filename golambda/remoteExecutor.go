package golambda

func invokeGrpcRequest(targetIp, targetBootstrapBinary string, targetPort uint16, event Event) (response Response, err error) {
	log.Debugf("Remote gRPC invocation to [%s:%v] with bootstrap [%s]", targetIp, targetPort, targetBootstrapBinary)
}
