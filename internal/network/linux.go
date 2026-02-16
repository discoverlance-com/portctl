//go:build linux

package network

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type LinuxManager struct{}

func (l LinuxManager) ListListeningProcesses() ([]LocalProcess, error) {
	cmd := exec.Command("ss", "-ltnp")

	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(output), "\n")

	var results []LocalProcess

	var pidRegex = regexp.MustCompile(`pid=(\d+)`)

	for _, line := range lines {
		if !strings.Contains(line, "LISTEN") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 5 {
			continue
		}

		address := fields[3] // 0.0.0.0:8080 - local address
		partsIndex := strings.LastIndex(address, ":")
		portStr := address[partsIndex+1:]

		port, err := strconv.Atoi(portStr)
		if err != nil {
			continue
		}

		// Extract PID from users:(("tour",pid=1234,fd=6))
		matches := pidRegex.FindStringSubmatch(line)
		if len(matches) != 2 {
			continue
		}

		pid, err := strconv.Atoi(matches[1])
		if err != nil {
			continue
		}

		results = append(results, LocalProcess{
			LocalPort:     port,
			OwningProcess: pid,
		})
	}

	return results, nil
}

func (l LinuxManager) KillProcess(pid int) error {
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
	return LinuxManager{}
}
