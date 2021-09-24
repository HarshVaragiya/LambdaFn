proto:
	protoc -I=liblambda --go_out=plugins=grpc:. liblambda/grpc-api.proto
runtime:
	cd lambdaRuntime && go build -o ../bin/runtime .
base_images:
	docker build -t amazonlinux-python -f images/python.dockerfile .