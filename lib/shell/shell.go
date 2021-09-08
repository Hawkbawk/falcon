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

// RunCommands runs the specified list of commands through the bash shell. This means you can simply
// write your commands like "echo "hello world" | sudo tee /bin/useless > /dev/null". If any errors
// are encountered while running the script, this function returns them. Note that because this will
// run whatever commands you provide to it, you should not run user-provided input through it
// without at least sanitizing it.
func RunCommand(command string) error {
	cmd := exec.Command("bash", "-c", command)

	result, err := cmd.Output()

	if err != nil {
		return fmt.Errorf("Command(s) failed to run due to the following error: %v", result)
	}
	return nil
}
