package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
	//	"io/ioutil"
	"bytes"
	"io/ioutil"
	"log"
)

type OptionalString struct {
	ptr *string
}

func (s OptionalString) String() string {
	if s.ptr == nil {
		return ""
	}
	return *s.ptr
}

func Env(key string) OptionalString {
	value, ok := os.LookupEnv(key)
	if !ok {
		return OptionalString{nil}
	}
	return OptionalString{&value}
}

func Default(args ...interface{}) (string, error) {
	for _, arg := range args {
		if arg == nil {
			continue
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
		default:
			return "", fmt.Errorf("Default: unsupported type '%T'!", v)
		}
	}

	return "", errors.New("Default: all arguments are nil!")
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

var funcMap = template.FuncMap{
	"env":     Env,
	"default": Default,
	"require": Require,
}

func generateTemplate(source, name string) (string, error) {
	var t *template.Template
	var err error
	if t, err = template.New(name).Funcs(funcMap).Parse(source); err != nil {
		return "", err
	}
	var buffer bytes.Buffer
	if err = t.Execute(&buffer, nil); err != nil {
		return "", err
	}
	return buffer.String(), nil
}

func generateFile(templatePath, destinationPath string, debugTemplates bool) error {
	if !filepath.IsAbs(templatePath) {
		return fmt.Errorf("Template path '%s' is not absolute!", templatePath)
	}

	if !filepath.IsAbs(destinationPath) {
		return fmt.Errorf("Destination path '%s' is not absolute!", destinationPath)
	}

	var slice []byte
	var err error
	if slice, err = ioutil.ReadFile(templatePath); err != nil {
		return err
	}
	s := string(slice)
	result, err := generateTemplate(s, filepath.Base(templatePath))
	if err != nil {
		return err
	}

	if debugTemplates {
		log.Printf("Printing parsed template to stdout. (It's delimited by 2 character sequence of '\\x00\\n'.)\n%s\x00\n", result)
	}

	if err = ioutil.WriteFile(destinationPath, []byte(result), 0664); err != nil {
		return err
	}

	return nil
}
