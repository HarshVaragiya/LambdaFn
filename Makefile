proto:
	protoc -I=liblambda --go_out=plugins=grpc:. liblambda/grpc-api.proto
runtime:
	cd lambdaRuntime && go build -o ../bin/runtime .
images: runtime
	docker build -t amazonlinux-python -f images/python/Dockerfile .
python_test_container: images
	docker run --rm --name python_test_container -p 8888:8888 -v $(shell pwd)/bin:/lambda/code -e LAMBDA_HANDLER_FUNCTION=example.lambda_handler -e LAMBDA_FUNCTION_NAME=my-awesome-function amazonlinux-python