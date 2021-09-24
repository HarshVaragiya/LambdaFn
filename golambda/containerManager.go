package golambda

func BuildLambdaImage(functionName,codeUri,runtime,handler string, environ map[string]string){
	log.Printf("Generating new function image for [%v] with runtime [%v]", functionName, runtime)
	
}
