package exec

import "os/exec"

type PowershellExecutor struct{}

func (r *PowershellExecutor) RunPowerShellCommand(command string) ([]byte, error) {
	cmd := exec.Command(
		"powershell",
		"-NoProfile",
		"-NonInteractive",
		"-Command",
		command,
	)

	output, err := cmd.CombinedOutput()

	if err != nil {
		return output, err
	}

	return output, nil
}
