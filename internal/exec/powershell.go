package exec

import "os/exec"

func RunPowerShellCommand(command string) ([]byte, error) {
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
