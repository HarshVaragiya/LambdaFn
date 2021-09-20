# LambdaFn
An open source, self-hosted service to run AWS Lambda Functions on your own hardware.

## Project
- This project is to try creating a service which can work like AWS Lambda / Google AppEngine
- It would lack the trigger integrations that AWS Provides with its lambda functions
- It would allow user to build a lambda function via AWS Cloudformation templates (drop in replacement)
- It would allow users to execute the lambda and send events via an HTTP REST API 

## Goals 
The service should allow executing the following:
- Binary on the server running the service
- Code written as Lambda zip running inside container on the server
- Container to be executed on the server (or kubernetes cluster)
- (practically any struct implementing the LambdaExecutor Interface can execute code)


