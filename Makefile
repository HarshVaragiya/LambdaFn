proto:
	protoc -I=liblambda --go_out=plugins=grpc:. liblambda/grpc-api.proto
runtime:
	cd lambdaRuntime && go build -o ../bin/runtime .
images: runtime
	docker build -t amazonlinux-python -f images/python/Dockerfile .
install: images
	go install .
