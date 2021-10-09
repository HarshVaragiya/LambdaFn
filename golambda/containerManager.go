package golambda

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

type DockerContainerManager struct {
	dockerClient *client.Client
}

func (manager *DockerContainerManager) Init() error {
	var err error
	manager.dockerClient, err = client.NewClientWithOpts(client.WithAPIVersionNegotiation())
	if err != nil {
		log.Errorf("error connecting to docker daemon. error = %v", err)
		return err
	}
	return nil
}

func (manager *DockerContainerManager) startContainer(ctx context.Context, image, localPort, srcFolder, targetFolder string, env map[string]string) (string, error) {
	containerOptions := &container.Config{
		Image:        image,
		ExposedPorts: nat.PortSet{nat.Port(fmt.Sprintf("%d", containerRpcPort)): struct{}{}},
		Env:          generateEnvStrings(env),
	}
	hostOptions := &container.HostConfig{
		PortBindings: map[nat.Port][]nat.PortBinding{
			nat.Port(fmt.Sprintf("%d", containerRpcPort)): {{HostIP: "127.0.0.1", HostPort: localPort}},
		},
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: srcFolder,
				Target: targetFolder,
			},
		},
		AutoRemove: true,
	}
	resp, err := manager.dockerClient.ContainerCreate(ctx, containerOptions, hostOptions, nil, nil, "")
	if err != nil {
		log.Errorf("error creating container. error = %v", err)
		return "", err
	}
	if err = manager.dockerClient.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		log.Errorf("error starting the container. error = %v", err)
		return "", err
	}
	log.Debugf("container [%s] started with required mappings", resp.ID)
	return resp.ID, nil
}

func (manager *DockerContainerManager) stopContainer(ctx context.Context, containerId string) error {
	err := manager.dockerClient.ContainerStop(ctx, containerId, &containerTimeout)
	if err != nil {
		log.Errorf("error stopping docker container [%s]. error = %v", containerId, err)
	}
	return err
}

func generateEnvStrings(env map[string]string) []string {
	var strings []string
	for key, value := range env {
		strings = append(strings, fmt.Sprintf("%s=%s", key, value))
	}
	return strings
}

func prepareEnvironmentVariables(functionName, handler string, env map[string]string) map[string]string {
	if env == nil {
		log.Debug("env is not defined")
		env = make(map[string]string)
	}
	_, exists := env["LAMBDA_FUNCTION_NAME"]
	if exists {
		log.Tracef("Lambda environment specifies function name. using that.")
	} else {
		env["LAMBDA_FUNCTION_NAME"] = functionName
	}

	existingHandler, exists := env["LAMBDA_HANDLER_FUNCTION"]
	if exists {
		if existingHandler != handler {
			log.Errorf("Lambda environment specifies handler different than the one configured. overriding")
			env["LAMBDA_HANDLER_FUNCTION"] = handler
		} else {
			log.Tracef("Lambda environment has pre-configured handler")
		}
	} else {
		env["LAMBDA_HANDLER_FUNCTION"] = handler
	}
	return env
}
