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

func TestRedactor(t *testing.T) {
	t.Run("leave alone commands that do not match redaction filter", func(t *testing.T) {
		testCases := []struct{ matcher, transformation, source, expected string }{
			{"bob", "sam", "ls /sam", "ls /sam"},
			{"sam", "bob", "command bob", "command bob"},
		}
		for _, testCase := range testCases {
			var redactor Transformer = Redactor{
				testCase.matcher: testCase.transformation,
			}
			t.Run(testCase.source, func(t *testing.T) {
				actual := redactor.transform(testCase.source)
				if actual != testCase.expected {
					t.Errorf("Expected %q, got %q", testCase.expected, actual)
				}
			})
		}
	})

	t.Run("modifies commands that match redaction filter", func(t *testing.T) {
		testCases := []struct{ matcher, transformation, source, expected string }{
			{`(--pass)=\w+`, "$1=REDACTED", "--pass=bob", "--pass=REDACTED"},
			{"sam", "bob", "do sam", "do bob"},
		}
		for _, testCase := range testCases {
			var redactor Transformer = Redactor{
				testCase.matcher: testCase.transformation,
			}
			t.Run(testCase.source, func(t *testing.T) {
				actual := redactor.transform(testCase.source)
				if actual != testCase.expected {
					t.Errorf("Expected %q, got %q", testCase.expected, actual)
				}
			})
		}
	})
}
