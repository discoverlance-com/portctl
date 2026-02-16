//go:build windows

package network

import (
	"encoding/json"
	"fmt"

	"github.com/discoverlance-com/portctl/internal/exec"
)

type WindowsManager struct{}

func (w WindowsManager) ListListeningProcesses() ([]LocalProcess, error) {
	script := `
		@(Get-NetTCPConnection -State Listen |
		Select-Object LocalPort, OwningProcess) |
		ConvertTo-Json -Depth 3
	`

	output, err := exec.RunPowerShellCommand(script)

	if err != nil {
		return nil, err
	}

	var result []LocalProcess

	if jsonErr := json.Unmarshal(output, &result); jsonErr != nil {
		return nil, jsonErr
	}

	seen := make(map[int]bool)
	var unique []LocalProcess

	for _, p := range result {
		if !seen[p.LocalPort] {
			seen[p.LocalPort] = true
			unique = append(unique, p)
		}
	}

	return unique, nil
}

func (l WindowsManager) KillProcess(pid int) error {
	command := fmt.Sprintf("Stop-Process -Id %d -Force -ErrorAction Stop; Write-Output 'OK'", pid)

	output, err := exec.RunPowerShellCommand(command)
	if err != nil {
		return fmt.Errorf("Failed to stop process: %v\nOutput: %s", err, output)
	}

	return nil
}

func NewManager() PortManager {
	return WindowsManager{}
}
