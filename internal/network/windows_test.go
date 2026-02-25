//go:build windows

package network

import (
	"errors"
	"testing"
)

type fakeExecutor struct {
	output []byte
	err    error
}

func (f *fakeExecutor) RunPowerShellCommand(script string) ([]byte, error) {
	return f.output, f.err
}

func TestListListeningProcesses(t *testing.T) {
	tests := []struct {
		name        string
		output      string
		executorErr error
		want        []LocalProcess
		wantErr     bool
	}{
		{
			"Valid output returns one process",
			`[{"LocalPort":8080,"OwningProcess":1234}]`,
			nil,
			[]LocalProcess{{LocalPort: 8080, OwningProcess: 1234}},
			false,
		},
		{
			"malformed JSON returns error",
			`I am free`,
			nil,
			nil,
			true,
		},
		{
			"duplicate ports are deduplicated",
			`[{"LocalPort":8080,"OwningProcess":1234},{"LocalPort":8080,"OwningProcess":1234}]`,
			nil,
			[]LocalProcess{{LocalPort: 8080, OwningProcess: 1234}},
			false,
		},
		{
			"Empty JSON returns an empty slice",
			`[]`,
			nil,
			[]LocalProcess{},
			false,
		},
		{
			"Command failure returns an error",
			``,
			errors.New("powershell command failed"),
			nil,
			true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			manager := NewManagerWithExecutor(&fakeExecutor{
				output: []byte(test.output),
				err:    test.executorErr,
			})

			result, err := manager.ListListeningProcesses()

			if (err != nil) != test.wantErr {
				t.Errorf("manager.ListListeningProcesses() error = %v, wantErr %v", err, test.wantErr)
				return
			}

			if test.wantErr {
				return
			}

			if len(result) != len(test.want) {
				t.Fatalf("manager.ListListeningProcesses() len = %d, want %d", len(result), len(test.want))
			}

			if len(result) > 0 {
				if result[0].LocalPort != test.want[0].LocalPort || result[0].OwningProcess != test.want[0].OwningProcess {
					t.Fatalf("manager.ListListeningProcesses()[0] = %+v, want %+v", result[0], test.want[0])
				}
			}
		})
	}

}

func TestKillProcess(t *testing.T) {
	tests := []struct {
		name        string
		pid         int
		output      string
		executorErr error
		wantErr     bool
	}{
		{
			"successfully kills process",
			1234,
			"OK",
			nil,
			false,
		},
		{
			"command failure returns error",
			1234,
			"Access denied",
			errors.New("exit code 1"),
			true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			manager := NewManagerWithExecutor(&fakeExecutor{
				output: []byte(test.output),
				err:    test.executorErr,
			})

			err := manager.KillProcess(test.pid)

			if (err != nil) != test.wantErr {
				t.Errorf("manager.KillProcess() error = %v, wantErr %v", err, test.wantErr)
			}
		})
	}
}
