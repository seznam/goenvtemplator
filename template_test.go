package main

import (
	"github.com/joho/godotenv"
	"os"
	"reflect"
	"testing"
	"text/template"
)

func TestGenerateTemplate(t *testing.T) {
	templateName := "test"
	var tests = []struct {
		in   string
		want string
		err  error
	}{
		{`K={{ env "GOENVTEMPLATOR_DEFINED_VAR" }}`, `K=foo`, nil},
		{`K={{ env "GOENVTEMPLATOR_DEFINED_FILE_VAR" }}`, `K=bar`, nil},
		{`K={{ env "NONEXISTING" }}`, `K=`, nil},
		{`K={{ .NONEXISTING }}`, ``, template.ExecError{}},
		{`K={{ .NonExisting | default "default value" }}`, ``, template.ExecError{}},
		{`K={{ env "GOENVTEMPLATOR_DEFINED_VAR" | default "xxx" }}`, `K=foo`, nil},
		{`K={{ env "GOENVTEMPLATOR_DEFINED_FILE_VAR" | default "xxx" }}`, `K=bar`, nil},
		{`K={{ env "NONEXISTING"| default "default value" }}`, `K=default value`, nil},
		{`{{ "hi!" | upper | repeat 3 }}`, `HI!HI!HI!`, nil},
		{`{{$v := "foo/bar/baz" | split "/"}}{{$v._1}}`, `bar`, nil},
		{`<?xml version="1.0"?>`, `<?xml version="1.0"?>`, nil},
	}

	os.Setenv("GOENVTEMPLATOR_DEFINED_VAR", "foo")

	err := godotenv.Load("./tests/fixtures.env")
	if err != nil {
		t.Errorf("Cannot load env file: %q", err)
	}

	for _, testcase := range tests {
		got, gotErr := generateTemplate(testcase.in, templateName)

		if testcase.want != got {
			t.Errorf("generateTemplate(%q, %q) => (%q, _), want (%q, _)", testcase.in, templateName, got, testcase.want)
		}

		errType, gotErrType := reflect.TypeOf(testcase.err), reflect.TypeOf(gotErr)

		if errType != gotErrType {
			t.Errorf("generateTemplate(%q, %q) => (_, %q), want (_, %q)", testcase.in, templateName, gotErrType, errType)
		}
	}
}
