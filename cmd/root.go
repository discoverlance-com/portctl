package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/discoverlance-com/portctl/internal/network"
)

func Execute() {
	flag.Usage = func() {
		fmt.Printf("Usage: portctl <command> [arguments]\n\n")
		fmt.Println("Available commands:")
		fmt.Println("  list    Show all listening ports and their PIDs")
		fmt.Println("  kill    Terminate a process by port or PID")
		fmt.Println("\nUse \"portctl <command> -help\" for more information about a command.")
	}
	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	manager := network.NewManager()

	switch os.Args[1] {
	case "list":
		HandleListRunningProcessesCommand(manager)
	case "kill":
		HandleKillProcessCommand(manager)
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", os.Args[1])
		fmt.Fprintln(os.Stderr, "Available commands: list, kill")
		os.Exit(1)
	}
}
