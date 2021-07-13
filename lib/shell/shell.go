package shell

import (
	"errors"
	"fmt"
	"os/exec"
)

type CommandNotFound struct {
	message string
}

func (e CommandNotFound) Error() string {
	return e.message
}

func RunCommand(desiredCommand string, args []string) error {
	execCommand, err := exec.LookPath(desiredCommand)

	if err != nil {
		return CommandNotFound{message: fmt.Sprint("Couldn't find command", desiredCommand)}
	}

	executable := &exec.Cmd{
		Path: execCommand,
		Args: args,
	}

	if err = executable.Run(); err != nil {
		return errors.New(fmt.Sprintln("Command", desiredCommand, "failed with error:", err.Error()))
	}

	return nil
}