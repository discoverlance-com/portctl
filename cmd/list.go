package cmd

import (
	"fmt"
	"os"

	"github.com/discoverlance-com/portctl/internal/network"
)

func HandleListRunningProcessesCommand(manager network.PortManager) {
	runningProcesses, err := manager.ListListeningProcesses()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to list running processes: %v", err)
		os.Exit(1)
	}

	fmt.Printf("%-10s %-10s\n", "Port", "PID")

	for _, process := range runningProcesses {
		fmt.Printf("%-10d %-10d\n", process.LocalPort, process.OwningProcess)
	}
}
