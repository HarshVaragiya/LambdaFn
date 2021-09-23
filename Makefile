proto:
	protoc -I=liblambda --go_out=plugins=grpc:. liblambda/grpc-api.proto