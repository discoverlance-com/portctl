//go:build darwin

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
			"valid lsof output returns processes",
			"p1234\nn*:8080\np5678\nn127.0.0.1:3000\n",
			nil,
			[]LocalProcess{
				{LocalPort: 8080, OwningProcess: 1234},
				{LocalPort: 3000, OwningProcess: 5678},
			},
			false,
		},
		{
			"non-lsof output returns empty results",
			`I am free`,
			nil,
			[]LocalProcess{},
			false,
		},
		{
			"empty output returns an empty slice",
			``,
			nil,
			[]LocalProcess{},
			false,
		},
		{
			"executor failure returns an error",
			``,
			errors.New("lsof command failed"),
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

			for i := range result {
				if result[i].LocalPort != test.want[i].LocalPort || result[i].OwningProcess != test.want[i].OwningProcess {
					t.Fatalf("manager.ListListeningProcesses()[%d] = %+v, want %+v", i, result[i], test.want[i])
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
