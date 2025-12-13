package main

import (
	"errors"
	"os"
	"testing"
)

func Test_run(t *testing.T) {
	tests := []struct {
		name string   // description of this test case
		args []string // input parameters
		err  error
	}{
		{
			name: "no command provided",
			args: []string{os.Args[0]},
			err:  ErrMissingCommand,
		},
		{
			name: "unknow command",
			args: []string{os.Args[0], "patata"},
			err:  ErrUnknownCommand,
		},
		{
			name: "everything is fine (no error)",
			args: []string{os.Args[0], "start"},
			err:  nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := run(tt.args); !errors.Is(err, tt.err) {
				t.Errorf("parseArgs() - %s test failed\n  wanted %s, got %s\n", tt.name, tt.err, err)
			}
		})
	}
}

func Test_start(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		args    []string
		wantErr error
	}{
		{
			name:    "placeholder",
			args:    []string{},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := start(tt.args)
			if !errors.Is(tt.wantErr, gotErr) {
				t.Errorf("start() - test %s failed: wanted %s, got %v", tt.name, tt.wantErr, gotErr)
			}
		})
	}
}
