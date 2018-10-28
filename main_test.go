package main

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"

	spb "github.com/ebastos/shell-history/history"
)

type mockRedactor struct {
	input  []string
	output []string
}

func (redactor *mockRedactor) transform(source []string) (output []string) {
	redactor.input = source
	return redactor.output
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
		redactor := &mockRedactor{output: tt.args.argsWithoutProg}
		t.Run(tt.name, func(t *testing.T) {
			got := getinformation(
				redactor, tt.args.argsWithoutProg, tt.args.commandExitCode)
			if len(got.Command) != len(tt.args.argsWithoutProg) {
				t.Errorf("getinformation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCommandRedactions(t *testing.T) {
	t.Run("parses configuration file definitions", func(t *testing.T) {
		redactor := Redactor{"thing": "value"}
		configFile := bytes.NewBufferString(`{"redactors":{"thing": "value"}}`)
		config := initConfig(configFile)
		if !reflect.DeepEqual(config.Redactors, redactor) {
			t.Errorf("Expected config.Redactors to equal %q, not %q",
				redactor, config.Redactors)
		}
	})

	t.Run("information retrieval triggers redactor logic", func(t *testing.T) {
		redactor := &mockRedactor{output: []string{"mocked", "output"}}
		commandArguments := []string{"some", "command"}
		command := getinformation(redactor, commandArguments, 0)
		if !reflect.DeepEqual(redactor.input, commandArguments) {
			t.Errorf("Expected redactor to be passed %q, but was passed %q",
				commandArguments, redactor.input)
		}
		if !reflect.DeepEqual(command.Command, redactor.output) {
			t.Errorf("Expected getinformation.Command to return %q, not %q",
				redactor.output, command.Command)
		}
	})

	t.Run("leave alone commands that don't match redactor", func(t *testing.T) {
		testCases := []struct {
			match, transform string
			source, expected []string
		}{
			{"cyclops", "scott", []string{"ls", "scott"}, []string{"ls", "scott"}},
			{"storm", "ororo", []string{"cmd", "ororo"}, []string{"cmd", "ororo"}},
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
