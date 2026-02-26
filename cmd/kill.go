package cmd

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/discoverlance-com/portctl/internal/network"
)

func HandleKillProcessCommand(manager network.PortManager) {
	killCmd := flag.NewFlagSet("kill", flag.ExitOnError)

	port := killCmd.Int("port", 0, "The Port of the running service you want to kill")
	process := killCmd.Int("pid", 0, "The Process ID for the running service you want to kill")
	confirm := killCmd.Bool("y", false, "Confirmation that you want to kill the process")

	err := killCmd.Parse(os.Args[2:])

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing command args: %v\n", err)
		killCmd.Usage()
		os.Exit(1)
	}

	err = ValidateKillProcessFlags(*port, *process)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		killCmd.Usage()
		os.Exit(1)
	}

	processes, err := manager.ListListeningProcesses()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to list running processes: %v\n", err)
		os.Exit(1)
	}

	var processToKill int

	for index := range processes {
		if processes[index].LocalPort == *port {
			processToKill = processes[index].OwningProcess
			break
		}

		if processes[index].OwningProcess == *process {
			processToKill = processes[index].OwningProcess
			break
		}
	}

	if processToKill == 0 {
		fmt.Fprintln(os.Stderr, "Process not found")
		os.Exit(1)
	}

	// confirm
	if !*confirm {
		fmt.Print("Are you sure you want to continue? Y/N: ")

		var confirmation string
		_, err := fmt.Scanln(&confirmation)

		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to confirm your action, please type Y/N to confirm.")
			os.Exit(1)
		}

		if strings.ToLower(confirmation) != "y" {
			fmt.Println("Cancelled command")
			os.Exit(0)
		}
	}

	err = manager.KillProcess(processToKill)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to kill running process: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Process: %d killed successfully\n", processToKill)
}

func ValidateKillProcessFlags(port int, process int) error {
	if port == 0 && process == 0 {
		return errors.New("provide at least one flag, -process or -port")
	}

	if port != 0 && process != 0 {
		return errors.New("provide only one of the flags, -process or -port")
	}

	return nil
}
