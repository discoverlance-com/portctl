//go:build darwin

package network

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type DarwinManager struct{}

func (d DarwinManager) ListListeningProcesses() ([]LocalProcess, error) {
	cmd := exec.Command(
		"lsof",
		"-n",
		"-P",
		"-iTCP",
		"-sTCP:LISTEN",
		"-F", "pn",
	)

	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(output), "\n")

	var results []LocalProcess
	var currentPID int

	for _, line := range lines {
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
		return fmt.Errorf("Could not find process %d: %w", pid, err)
	}

	if err := proc.Kill(); err != nil {
		return fmt.Errorf("Failed to kill process %d: %w", pid, err)
	}

	return nil
}

func NewManager() PortManager {
	return DarwinManager{}
}
