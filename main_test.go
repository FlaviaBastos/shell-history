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
	t.Run("leave alone commands that don't match redactor", func(t *testing.T) {
		testCases := []struct{ match, transform, source, expected string }{
			{"cyclops", "scott", "ls /scott", "ls /scott"},
			{"storm", "ororo", "command ororo", "command ororo"},
		}
		for _, testCase := range testCases {
			var redactor Transformer = Redactor{
				testCase.match: testCase.transform,
			}
			t.Run(testCase.source, func(t *testing.T) {
				actual := redactor.transform(testCase.source)
				if actual != testCase.expected {
					t.Errorf("Expected %q, got %q", testCase.expected, actual)
				}
			})
		}
	})

	t.Run("modifies commands that match redactor", func(t *testing.T) {
		testCases := []struct{ match, transform, source, expected string }{
			{`(--pass)=\w+`, "$1=REDACTED", "--pass=peter", "--pass=REDACTED"},
			{"peter parker", "spider-man", "do peter parker", "do spider-man"},
		}
		for _, testCase := range testCases {
			var redactor Transformer = Redactor{
				testCase.match: testCase.transform,
			}
			t.Run(testCase.source, func(t *testing.T) {
				actual := redactor.transform(testCase.source)
				if actual != testCase.expected {
					t.Errorf("Expected %q, got %q", testCase.expected, actual)
				}
			})
		}
	})

	t.Run("modifies commands that match multiple redactors", func(t *testing.T) {
		var redactor Transformer = Redactor{
			`(--pass)=\w+`: "$1=REDACTED",
			`(--key)=\w+`:  "$1=NOWAY",
		}
		expected := "cmd --pass=REDACTED --key=NOWAY"
		actual := redactor.transform("cmd --pass=clark --key=kent")
		if actual != expected {
			t.Errorf("Expected %q, got %q", expected, actual)
		}
	})
}
