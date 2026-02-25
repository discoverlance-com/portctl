//go:build linux

package network

import (
	"errors"
	"os/exec"
	"testing"
)

type fakeExecutor struct {
	output []byte
	err    error
}

func (f *fakeExecutor) RunLinuxCommand(command string, args ...string) ([]byte, error) {
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
			`LISTEN 0 4096 0.0.0.0:8080 0.0.0.0:* users:(("app",pid=1234,fd=6))`,
			nil,
			[]LocalProcess{{LocalPort: 8080, OwningProcess: 1234}},
			false,
		},
		{
			"malformed output should return no error + empty results",
			`I am free`,
			nil,
			[]LocalProcess{},
			false,
		},
		{
			"Empty results returns an empty slice",
			``,
			nil,
			[]LocalProcess{},
			false,
		},
		{
			"Command executor failure returns an error",
			``,
			errors.New("linux, ss command failed"),
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

func TestKillProcess_Integration(t *testing.T) {
	cmd := exec.Command("sleep", "30")
	if err := cmd.Start(); err != nil {
		t.Skipf("skipping integration test: unable to start sleep process: %v", err)
	}

	defer func() {
		if cmd.ProcessState == nil || !cmd.ProcessState.Exited() {
			_ = cmd.Process.Kill()
			_, _ = cmd.Process.Wait()
		}
	}()

	manager := NewManagerWithExecutor(&fakeExecutor{})
	if err := manager.KillProcess(cmd.Process.Pid); err != nil {
		t.Fatalf("manager.KillProcess() unexpected error: %v", err)
	}

	if err := cmd.Wait(); err == nil {
		t.Fatalf("expected process to be terminated by KillProcess, but wait returned nil")
	}
}
