package lambda

import (
	"bytes"
	"context"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io"
	"os/exec"
)

func RunLocalBinary(stdinInput string, ctx context.Context, codeUri, functionHandler string) (string, string, error) {
	cmd := exec.CommandContext(ctx, codeUri, functionHandler)
	log.Debugf("Executing local binary: '%s', with args: '%s' ", codeUri, functionHandler)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Warn("Unable to Open STDIN.")
		log.Warn(err)
		return "", "", err
	}
	_, err = io.WriteString(stdin, stdinInput)
	if err != nil {
		log.Warn("Unable to write to STDIN.")
		log.Warn(err)
		return "", "", err
	}
	var stdoutBytes bytes.Buffer
	var stderrBytes bytes.Buffer
	cmd.Stdout = &stdoutBytes
	cmd.Stderr = &stderrBytes
	err = cmd.Run()
	if err != nil {
		log.Warn("Error running the required binary")
		log.Warn(err)
	}
	return stdoutBytes.String(), stderrBytes.String(), err
}

func EventToStdinInput(event *Event) (string, error) {
	eventBytes, err := json.Marshal(event)
	if err != nil{
		log.Warnf("Error Marshalling Event JSON. EventId : %s", event.EventId)
		return "", err
	} else {
		return string(eventBytes), nil
	}
}

