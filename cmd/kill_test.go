package cmd

import "testing"

func TestValidateKillProcessFlags(t *testing.T) {
	tests := []struct {
		name    string
		port    int
		process int
		wantErr bool
	}{
		{"Both zero", 0, 0, true},
		{"Both set", 8080, 1234, true},
		{"Only port", 8080, 0, false},
		{"Only process", 0, 1234, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := ValidateKillProcessFlags(test.port, test.process)
			if (err != nil) != test.wantErr {
				t.Errorf("ValidateKillProcessFlags() error = %v, wantErr %v", err, test.wantErr)
			}
		})
	}
}
