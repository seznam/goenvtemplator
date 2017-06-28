package main

import (
	"github.com/joho/godotenv"
	"os"
	"testing"
)

func TestGenerateTemplate(t *testing.T) {
	var tests = []struct {
		in   string
		want string
		err  error
	}{
		{`K={{ env "GOENVTEMPLATOR_DEFINED_VAR" }}`, `K=foo`, nil},
		{`K={{ env "GOENVTEMPLATOR_DEFINED_FILE_VAR" }}`, `K=bar`, nil},
		{`K={{ env "NONEXISTING" }}`, `K=`, nil},
		{`K={{ .NONEXISTING }}`, `K=`, nil},
		{`K={{ .NonExisting | default "default value" }}`, `K=default value`, nil},
		{`K={{ env "GOENVTEMPLATOR_DEFINED_VAR" | default "xxx" }}`, `K=foo`, nil},
		{`K={{ env "GOENVTEMPLATOR_DEFINED_FILE_VAR" | default "xxx" }}`, `K=bar`, nil},
		{`K={{ env "NONEXISTING"| default "default value" }}`, `K=default value`, nil},
		{`{{ "hi!" | upper | repeat 3 }}`, `HI!HI!HI!`, nil},
		{`{{$v := "foo/bar/baz" | split "/"}}{{$v._1}}`, `bar`, nil},
		{`<?xml version="1.0"?>`, `<?xml version="1.0"?>`, nil},
	}

	templateName := "test"

	os.Setenv("GOENVTEMPLATOR_DEFINED_VAR", "foo")

	err := godotenv.Load("./tests/fixtures.env")
	if err != nil {
		t.Errorf("Cannot load env file: %q", err)
	}

	for _, tt := range tests {
		got, gotErr := generateTemplate(tt.in, templateName)

		if tt.want != got {
			t.Errorf("generateTemplate(%q, %q) => (%q, _), want (%q, _)", tt.in, templateName, got, tt.want)
		}

		if tt.err != gotErr {
			t.Errorf("generateTemplate(%q, %q) => (_, %q), want (_, %q)", tt.in, templateName, gotErr, tt.err)
		}
	}
}
