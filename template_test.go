package main

import (
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"text/template"
)

func absRelPath(t *testing.T, relPath string) (string, string) {
	absPath, err := filepath.Abs(relPath)
	if err != nil {
		t.Fatal(err)
	}
	return absPath, relPath
}

func TestGenerateFile(t *testing.T) {
	templatePathRel, templatePathAbs := absRelPath(t, "./tests/test.tmpl")
	outputPathRel, outputPathAbs := absRelPath(t, "./tests/test.out")

	testCases := []struct {
		name          string
		templatePath  string
		outputPath    string
		expectedError bool
	}{
		{name: "template absolute path, output absolute path", templatePath: templatePathAbs, outputPath: outputPathAbs, expectedError: false},
		{name: "template absolute path, output relative path", templatePath: templatePathAbs, outputPath: outputPathRel, expectedError: false},
		{name: "template relative path, output absolute path", templatePath: templatePathRel, outputPath: outputPathAbs, expectedError: false},
		{name: "template relative path, output relative path", templatePath: templatePathRel, outputPath: outputPathRel, expectedError: false},
		{name: "non existing template relative path", templatePath: "./fooo", outputPath: outputPathRel, expectedError: true},
		{name: "non existing template absolute path", templatePath: "/xxx/yyy/foo/bar", outputPath: outputPathRel, expectedError: true},
		{name: "creepy relative template path", templatePath: "./tests/../tests/test.tmpl", outputPath: outputPathRel, expectedError: false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := generateFile(tc.templatePath, tc.outputPath, false, "{{", "}}")
			if tc.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}

	// Cleanup the output test file
	_ = os.Remove(outputPathRel)
}

func TestGenerateTemplate(t *testing.T) {
	templateName := "test"
	var tests = []struct {
		in         string
		want       string
		contains   string
		err        error
		leftDelim  string
		rightDelim string
	}{
		{in: `K={{ env "GOENVTEMPLATOR_DEFINED_VAR" }}`, want: `K=foo`},
		{in: `K={{ env "GOENVTEMPLATOR_DEFINED_FILE_VAR" }}`, want: `K=bar`},
		{in: `K={{ env "NONEXISTING" }}`, want: `K=`},
		{in: `K={{ .NONEXISTING }}`, want: ``, err: template.ExecError{}},
		{in: `K={{ .NonExisting | default "default value" }}`, want: ``, err: template.ExecError{}},
		{in: `K={{ env "GOENVTEMPLATOR_DEFINED_VAR" | default "xxx" }}`, want: `K=foo`},
		{in: `K={{ env "GOENVTEMPLATOR_DEFINED_FILE_VAR" | default "xxx" }}`, want: `K=bar`},
		{in: `K={{ env "NONEXISTING"| default "default value" }}`, want: `K=default value`},
		{in: `{{ range $key, $value := envall }} {{ $key }}={{ $value }};{{ end }}`, contains: ` GOENVTEMPLATOR_DEFINED_VAR=foo;`},
		{in: `{{ "hi!" | upper | repeat 3 }}`, want: `HI!HI!HI!`},
		{in: `{{$v := "foo/bar/baz" | split "/"}}{{$v._1}}`, want: `bar`},
		{in: `<?xml version="1.0"?>`, want: `<?xml version="1.0"?>`},
		{in: `K={{env "GOENVTEMPLATOR_DEFINED_VAR"}}`, want: `K={{env "GOENVTEMPLATOR_DEFINED_VAR"}}`, err: nil, leftDelim: "[[", rightDelim: "]]"},
		{in: `K=[[env "GOENVTEMPLATOR_DEFINED_VAR"]]`, want: `K=foo`, err: nil, leftDelim: "[[", rightDelim: "]]"},
		{in: `K={{ require (env "FOO" )}}`, err: template.ExecError{}},
		{in: `K={{ require (env "GOENVTEMPLATOR_DEFINED_VAR" )}}`, want: `K=foo`},
		{in: `K={{ require (env "GOENVTEMPLATOR_DEFINED_VAR_EMPTY" )}}`, err: template.ExecError{}},
		{in: `K={{ required "message" (env "GOENVTEMPLATOR_DEFINED_VAR_EMPTY") }}`, err: template.ExecError{}},
		{in: `K={{ required "message" (env "GOENVTEMPLATOR_DEFINED_VAR_EMPTY" | default "foo") }}`, want: `K=foo`},
		{in: `K={{ required "message" "foo" }}`, want: `K=foo`},
		{in: `K={{ required "message" "" }}`, err: template.ExecError{}},
	}

	_ = os.Setenv("GOENVTEMPLATOR_DEFINED_VAR", "foo")
	_ = os.Setenv("GOENVTEMPLATOR_DEFINED_VAR_EMPTY", "")

	err := godotenv.Load("./tests/fixtures.env")
	if err != nil {
		t.Errorf("Cannot load env file: %q", err)
	}

	for _, testcase := range tests {
		got, gotErr := generateTemplate(testcase.in, templateName, testcase.leftDelim, testcase.rightDelim)

		if testcase.contains != "" {
			if !strings.Contains(got, testcase.contains) {
				t.Errorf("generateTemplate(%q, %q, %q, %q) => (%q, _), want containing (%q, _)", testcase.in, templateName, testcase.leftDelim, testcase.rightDelim, got, testcase.contains)
			}
		} else {
			if testcase.contains == "" && testcase.want != got {
				t.Errorf("generateTemplate(%q, %q, %q, %q) => (%q, _), want (%q, _)", testcase.in, templateName, testcase.leftDelim, testcase.rightDelim, got, testcase.want)
			}
		}

		errType, gotErrType := reflect.TypeOf(testcase.err), reflect.TypeOf(gotErr)

		if errType != gotErrType {
			t.Errorf("generateTemplate(%q, %q, %q, %q)) => (_, %q), want (_, %q)", testcase.in, templateName, testcase.leftDelim, testcase.rightDelim, gotErrType, errType)
		}
	}
}
