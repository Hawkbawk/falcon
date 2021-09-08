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

// RunCommands runs the specified list of commands through the bash shell. This means you can simply
// write your commands like "echo "hello world" | sudo tee /bin/useless > /dev/null". If any errors
// are encountered while running the script, this function returns them. Note that because this will
// run whatever commands you provide to it, you should not run user-provided input through it
// without at least sanitizing it.
func RunCommands(commands string) error {
	cmd := exec.Command("bash", "-c", commands)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Command(s) failed to run due to the following error: %v", err)
	}
	return nil
}
