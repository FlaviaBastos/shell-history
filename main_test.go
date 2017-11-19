package main

import (
	"testing"

	spb "github.com/ebastos/shell-history/history"
)

func Test_getinformation(t *testing.T) {
	type args struct {
		argsWithoutProg []string
		commandExitCode int64
	}
	tests := []struct {
		name string
		args args
		want spb.Command
	}{
		{name: "command without parameters", args: args{
			argsWithoutProg: []string{"ls"},
			commandExitCode: 0,
		}},
		{name: "command with two parameters", args: args{
			argsWithoutProg: []string{"ls", "-lha"},
			commandExitCode: 0,
		}},
		{name: "command with multiple parameters", args: args{
			argsWithoutProg: []string{"ps", "aux", "|", "grep", "test"},
			commandExitCode: 0,
		}},
		{name: "missing command", args: args{
			argsWithoutProg: []string{},
			commandExitCode: 0,
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getinformation(tt.args.argsWithoutProg, tt.args.commandExitCode)
			if len(got.Command) != len(tt.args.argsWithoutProg) {
				t.Errorf("getinformation() = %v, want %v", got, tt.want)
			}
		})
	}
}
