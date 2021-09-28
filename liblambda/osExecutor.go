package liblambda

import (
	"bytes"
	"context"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"os/exec"
)

func RunLocalBinary(stdinInput string, ctx context.Context, codeUri, functionHandler string) (string, string, string, error) {
	cmd := exec.CommandContext(ctx, codeUri, functionHandler)
	log.Debugf("Executing local binary: '%s', with args: '%s' ", codeUri, functionHandler)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Warn("Unable to Open STDIN.")
		log.Warn(err)
		return "", "","", err
	}
	if n, err := io.WriteString(stdin, stdinInput); err != nil {
		log.Warn("Unable to write to STDIN.")
		log.Warn(err)
		return "", "","", err
	} else {
		log.Printf("wrote %v bytes as event string", n)
	}
	err = stdin.Close()
	if err != nil{
		log.Warnf("unable to close stdin. this might cause the program to hault. continuing...")
	}
	var stdoutBytes bytes.Buffer
	var stderrBytes bytes.Buffer
	var responseBytes bytes.Buffer
	cmd.Stdout = &stdoutBytes
	cmd.Stderr = &stderrBytes
	responseReader, responseWriter, err := os.Pipe()
	cmd.ExtraFiles = []*os.File{responseWriter}

	err = cmd.Run()
	if err != nil {
		log.Warn("Error running the required binary")
		log.Warn(err)
		return stdoutBytes.String(), stderrBytes.String(), "", err
	}

	if err = responseWriter.Close(); err != nil{
		log.Errorf("error closing FD 3 writer. quitting. error = %v", err)
		return stdoutBytes.String(), stderrBytes.String(), "", err
	}
	n, err := responseBytes.ReadFrom(responseReader)
	if err != nil {
		log.Warnf("error reading FD 3 pipe. error = %v", err)
	} else {
		log.Debugf("got %v bytes response ", n)
	}

	if err = responseReader.Close(); err != nil {
		log.Errorf("error closing FD 3 reader. SKIPPING. error = %v", err)
	}

	return stdoutBytes.String(), stderrBytes.String(),responseBytes.String(), err
}

func EventToStdinInput(event *Event) (string, error) {
	eventBytes, err := json.Marshal(event)
	if err != nil{
		log.Warnf("error marshalling event JSON. eventId : %s", event.EventId)
		return "", err
	} else {
		return string(eventBytes), nil
	}
}

