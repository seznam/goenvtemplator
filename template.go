package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/Masterminds/sprig"
	"os"
	"strings"
	"text/template"
)

type TextTemplar struct {
	source     string
	name       string
	delimLeft  string
	delimRight string
}

type OptionalString struct {
	ptr *string
}

func (s OptionalString) String() string {
	if s.ptr == nil {
		return ""
	}
	return *s.ptr
}

func Require(arg interface{}) (string, error) {
	if arg == nil {
		return "", errors.New("Required argument is missing!")
	}

	switch v := arg.(type) {
	case string:
		return v, nil
	case *string:
		if v != nil {
			return *v, nil
		}
	case OptionalString:
		if v.ptr != nil {
			return *v.ptr, nil
		}
	}

	return "", fmt.Errorf("Requires: unsupported type '%T'!", v)
}

func EnvAll() (map[string]string, error) {
	res := make(map[string]string)

	for _, item := range os.Environ() {
		split := strings.Split(item, "=")
		res[split[0]] = strings.Join(split[1:], "=")
	}

	return res, nil
}

var funcMap = template.FuncMap{
	"require": Require,
	"envall":  EnvAll,
}

func (templar *TextTemplar) generateTemplate() (string, error) {
	t, err := template.New(templar.name).
		Delims(templar.delimLeft, templar.delimRight).
		Option("missingkey=error").
		Funcs(funcMap).
		Funcs(sprig.TxtFuncMap()).
		Parse(templar.source)

	if err != nil {
		return "", err
	}

	var buffer bytes.Buffer
	// hacking because go 1.7 fails to throw error, see https://github.com/golang/go/commit/277bcbbdcd26f2d64493e596238e34b47782f98e
	emptyHash := map[string]interface{}{}
	if err = t.Execute(&buffer, &emptyHash); err != nil {
		return "", err
	}

	return buffer.String(), nil
}
