package main

import (
	"fmt"
	"reflect"
	"testing"

	spb "github.com/ebastos/shell-history/history"
)

// Used for mocking out the Redactor dependency.
type mockRedactor struct {
	source []string
}

func (redactor *mockRedactor) transform(source []string) (output []string) {
	redactor.source = source
	return source
}

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

func TestCommandRedactions(t *testing.T) {
	t.Run("leave alone commands that don't match redactor", func(t *testing.T) {
		testCases := []struct {
			match, transform string
			source, expected []string
		}{
			{"cyclops", "scott", []string{"ls", "/scott"}, []string{"ls", "/scott"}},
			{"storm", "ororo", []string{"command", "ororo"}, []string{"command", "ororo"}},
		}
		for _, testCase := range testCases {
			var redactor Transformer = Redactor{
				testCase.match: testCase.transform,
			}
			t.Run(fmt.Sprint(testCase.source), func(t *testing.T) {
				actual := redactor.transform(testCase.source)
				if !reflect.DeepEqual(actual, testCase.expected) {
					t.Errorf("Expected %q, got %q", testCase.expected, actual)
				}
			})
		}
	})

	t.Run("modifies commands that match redactor", func(t *testing.T) {
		testCases := []struct {
			match, transform string
			source, expected []string
		}{
			{`(--pass)=\w+`, "$1=REDACTED", []string{"--pass=peter"}, []string{"--pass=REDACTED"}},
			{"peter parker", "spider-man", []string{"do", "peter parker"}, []string{"do", "spider-man"}},
		}
		for _, testCase := range testCases {
			var redactor Transformer = Redactor{
				testCase.match: testCase.transform,
			}
			t.Run(fmt.Sprint(testCase.source), func(t *testing.T) {
				actual := redactor.transform(testCase.source)
				if !reflect.DeepEqual(actual, testCase.expected) {
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
		expected := []string{"cmd", "--pass=REDACTED", "--key=NOWAY"}
		actual := redactor.transform([]string{"cmd", "--pass=clark", "--key=kent"})
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("Expected %q, got %q", expected, actual)
		}
	})
}
