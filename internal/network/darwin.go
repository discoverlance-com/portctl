//go:build darwin

package network

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	internalexec "github.com/discoverlance-com/portctl/internal/exec"
)

type DarwinManager struct {
	executor LinuxCommandExecutor
}

func (d DarwinManager) ListListeningProcesses() ([]LocalProcess, error) {
	output, err := d.executor.RunLinuxCommand(
		"lsof",
		"-n",
		"-P",
		"-iTCP",
		"-sTCP:LISTEN",
		"-F", "pn",
	)

	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(output), "\n")

	var results []LocalProcess
	var currentPID int

	for _, line := range lines {
		if len(line) == 0 {
			continue
		}

		switch line[0] {
		case 'p':
			pid, err := strconv.Atoi(line[1:])

			if err == nil {
				currentPID = pid
			}
		case 'n':
			// Name line (*:8080)
			parts := strings.Split(line[1:], ":")
			portStr := parts[len(parts)-1]

			port, err := strconv.Atoi(portStr)
			if err != nil {
				continue
			}

			results = append(results, LocalProcess{
				LocalPort:     port,
				OwningProcess: currentPID,
			})
		}
	}

	return results, nil
}

func (l DarwinManager) KillProcess(pid int) error {
	proc, err := os.FindProcess(pid)

	if err != nil {
		return fmt.Errorf("could not find process %d: %w", pid, err)
	}

	if err := proc.Kill(); err != nil {
		return fmt.Errorf("failed to kill process %d: %w", pid, err)
	}

	return nil
}

func NewManagerWithExecutor(executor LinuxCommandExecutor) DarwinManager {
	return DarwinManager{executor: executor}
}

func NewManager() PortManager {
	return NewManagerWithExecutor(&internalexec.LinuxShellExecutor{})
}
