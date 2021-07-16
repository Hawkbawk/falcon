package shell

import (
	"fmt"
	"os/exec"
)

// Indicates that a command was unable to be found on the machine's PATH.
type CommandNotFound struct {
	message string
}

// Indicates that a command failed to run.
type CommandFailed struct {
	message string
}

func (e CommandNotFound) Error() string {
	return e.message
}

func (e CommandFailed) Error() string {
	return e.message
}

// RunCommand runs the specified command with the specified arguments. If it can't find the desired
// command, it returns a CommandNotFound error. If the command fails, it returns a CommandFailed
// error.
func RunCommand(desiredCommand string, args []string, runAsSudo bool) error {
	execCommand, err := exec.LookPath(desiredCommand)

	if err != nil {
		return CommandNotFound{message: fmt.Sprint("Couldn't find command", desiredCommand)}
	}

	executable := &exec.Cmd{
		Path: execCommand,
		Args: args,
	}

	if err = executable.Run(); err != nil {
		return CommandFailed{message: fmt.Sprintln("Command", desiredCommand, "failed with error:", err.Error())}
	}

	return nil
}
