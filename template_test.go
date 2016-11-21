package main

import (
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
		{`K={{ env "NONEXISTING" }}`, `K=`, nil},
		{`K={{ .NONEXISTING }}`, `K=<no value>`, nil},
		{`K={{ default .NonExisting "default value" }}`, `K=default value`, nil},
		{`K={{ default (env "GOENVTEMPLATOR_DEFINED_VAR") }}`, `K=foo`, nil},
		{`K={{ default (env "NONEXISTING") "default value" }}`, `K=default value`, nil},
	}

	templateName := "test"

	os.Setenv("GOENVTEMPLATOR_DEFINED_VAR", "foo")

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
