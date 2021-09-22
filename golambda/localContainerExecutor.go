package golambda

type ContainerExecutor struct {
	functionName string
	runtime      string
	name         string
	rpcPort      uint16
	codeExecutor BasicCodeExecutor
}

func NewContainerExecutor(function Function) ContainerExecutor {
	executor := BasicCodeExecutor{codeUri: function.CodeUri, functionHandler: function.Handler, functionTimeout: defaultTimeout}
	containerExecutor := ContainerExecutor{codeExecutor: executor, runtime: function.Runtime, functionName: function.Name}
	return containerExecutor
}

func (executor ContainerExecutor) execute(event Event) (response Response, err error) {
	log.Debugf("Invoking Lambda [%s] in container.", executor.functionName)
}

func (executor ContainerExecutor) startContainer() (err error) {

}

func (executor ContainerExecutor) stopContainer() (err error) {

}
