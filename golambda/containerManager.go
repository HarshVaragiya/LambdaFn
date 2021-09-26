package golambda

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

type DockerImageBuilder struct {
	dockerClient		*client.Client
	dockerBuildOpts 	*types.ImageBuildOptions
}

func (builder *DockerImageBuilder) Init() error {
	var err error
	builder.dockerClient, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Errorf("error connecting to docker daemon. error = %v", err)
		return err
	}
	builder.dockerBuildOpts = &types.ImageBuildOptions{
		Dockerfile: "Dockerfile",
		Remove: true,
	}
	return nil
}

func (builder *DockerImageBuilder) BuildLambdaImage (ctx context.Context, functionName, codeUri, runtime, handler string, environ map[string]string) error {
	log.Printf("Generating new function image for [%v] with runtime [%v]", functionName, runtime)
	tempDir, err := os.MkdirTemp(os.TempDir(), fmt.Sprintf("LambdaFn-Build-%s-",functionName))
	if err != nil {
		log.Warnf("error creating temp file for docker image creation. error = %v", err)
		return err
	}
	defer os.RemoveAll(tempDir)
	log.Tracef("created temporary directory at %v", tempDir)
	templateFile, exists := runtimeToDockerfileMap[runtime]
	if !exists {
		log.Errorf("runtime [%v] has no docker template present. please consider defining it.", runtime)
		return fmt.Errorf("docker template not found")
	}
	log.Debugf("loading template for runtime [%v] - [%v]", runtime, templateFile)
	prepareEnvironmentVariables(functionName, handler, environ)
	err = createDockerfileWithEnv(templateFile, tempDir, environ)
	if err != nil {
		log.Errorf("error creating dockerfile from template. error = %v", err)
		return err
	}
	err = loadCodeZipFromCodeUri(tempDir, codeUri)
	if err != nil {
		log.Errorf("error with codeUri. unable to continue with build. error = %v", err)
		return err
	}
	tar, err := archive.TarWithOptions(tempDir, &archive.TarOptions{})
	if err != nil{
		log.Errorf("error creating docker context tar archive. error = %v", err)
		return err
	}
	defer tar.Close()
	res, err := builder.dockerClient.ImageBuild(ctx, tar, *builder.dockerBuildOpts)
	if err != nil {
		log.Warnf("error creating function image. error = %v", err)
		return err
	}
	defer res.Body.Close()
	err = print(res.Body)
	if err != nil {
		log.Printf("error building image? error = %v",err)
		return err
	}
	return nil
}
type ErrorLine struct {
	Error       string      `json:"error"`
	ErrorDetail ErrorDetail `json:"errorDetail"`
}

type ErrorDetail struct {
	Message string `json:"message"`
}

func print(rd io.Reader) error {
	var lastLine string
	scanner := bufio.NewScanner(rd)
	for scanner.Scan() {
		lastLine = scanner.Text()
		fmt.Println(scanner.Text())
	}
	errLine := &ErrorLine{}
	json.Unmarshal([]byte(lastLine), errLine)
	if errLine.Error != "" {
		return errors.New(errLine.Error)
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func createDockerfileWithEnv(templateFile, tempDir string, environ map[string]string)error{
	templateBytes, err := ioutil.ReadFile(templateFile)
	if err != nil {
		log.Errorf("error reading template file [%v]. error = %v", templateFile, err)
		return err
	}
	dockerfileString := string(templateBytes)
	envString := environmentVariablesToString(environ)
	strings.Replace( dockerfileString,"<<ENVIRONMENT_VARIABLES>>", envString, 1)
	err = ioutil.WriteFile(fmt.Sprintf("%s/Dockerfile",tempDir),[]byte(dockerfileString), 666)
	if err != nil{
		log.Errorf("error saving dockerfile at [%v]. error = %v", tempDir, err)
		return err
	}
	return nil
}

func loadCodeZipFromCodeUri(tempDir, codeUri string) error {
	// this function can change if we load from s3 or HDFS or some other storage service. ideal way would be to use interface
	codeBytes, err := ioutil.ReadFile(codeUri)
	if err != nil {
		log.Errorf("error reading from codeUri file. error = %v", err)
		return err
	}
	err = ioutil.WriteFile(fmt.Sprintf("%s/source_code.zip",tempDir), codeBytes, 0666)
	if err != nil {
		log.Errorf("error writing to build folder. error saving source_code.zip . error = %v", err)
		return err
	}
	return nil
}

func prepareEnvironmentVariables(functionName, handler string, env map[string]string){
	_ , exists := env["LAMBDA_FUNCTION_NAME"]
	if exists {
		log.Tracef("Lambda environment specifies function name. using that.")
	} else {
		env["LAMBDA_FUNCTION_NAME"] = functionName
	}

	existingHandler , exists := env["LAMBDA_HANDLER_FUNCTION"]
	if exists{
		if existingHandler != handler {
			log.Errorf("Lambda environment specifies handler different than the one configured. overriding")
			env["LAMBDA_HANDLER_FUNCTION"] = handler
		} else {
			log.Tracef("Lambda environment has pre-configured handler")
		}
	} else {
		env["LAMBDA_HANDLER_FUNCTION"] = handler
	}
}

func environmentVariablesToString(env map[string]string) string {
	envStrings := make([]string, len(env))
	for key, value := range env {
		log.Debugf("Adding key [%v] as [%v]", key, value)
		envStrings = append(envStrings, fmt.Sprintf("ENV %s=%s", key, value))
	}
	return strings.Join(envStrings, "\n") + "\n"
}