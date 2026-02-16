package network

type LocalProcess struct {
	LocalPort     int `json:"LocalPort"`
	OwningProcess int `json:"OwningProcess"`
}

type PortManager interface {
	ListListeningProcesses() ([]LocalProcess, error)
	KillProcess(int) error
}
