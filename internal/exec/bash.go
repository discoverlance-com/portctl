package exec

import "os/exec"

type LinuxShellExecutor struct{}

func (r *LinuxShellExecutor) RunLinuxCommand(name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)

	output, err := cmd.Output()

	if err != nil {
		return output, err
	}

	return output, nil
}
