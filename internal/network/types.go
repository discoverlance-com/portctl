package network

type LocalProcess struct {
	LocalPort     int `json:"LocalPort"`
	OwningProcess int `json:"OwningProcess"`
}

type PortManager interface {
	ListListeningProcesses() ([]LocalProcess, error)
	KillProcess(int) error
}

type WindowsCommandExecutor interface {
	RunPowerShellCommand(script string) ([]byte, error)
}

type LinuxCommandExecutor interface {
	RunLinuxCommand(name string, args ...string) ([]byte, error)
}
