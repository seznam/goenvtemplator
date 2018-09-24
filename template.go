package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/Masterminds/sprig"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func Require(arg string) (string, error) {
	fmt.Fprintf(os.Stderr, "WARNING: require built-in function is deprecated. Use required instead.\n")
	if len(arg) == 0 {
		return "", errors.New("Required argument is missing or empty!")
	}
	return arg, nil
}

// copied from Helm source:
// https://github.com/kubernetes/helm/blob/78d6b930bd325ed87b297c57b02fc7c9c7dfcfac/pkg/engine/engine.go#L156-L165
func Required(warn string, val interface{}) (interface{}, error) {
	if val == nil {
		return val, fmt.Errorf(warn)
	} else if _, ok := val.(string); ok {
		if val == "" {
			return val, fmt.Errorf(warn)
		}
	}
	return val, nil
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
	"required": Required,
}

func generateTemplate(source, name string, delimLeft string, delimRight string) (string, error) {
	var t *template.Template
	var err error
	t, err = template.New(name).Delims(delimLeft, delimRight).Option("missingkey=error").Funcs(funcMap).Funcs(sprig.TxtFuncMap()).Parse(source)
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

func generateFile(templatePath, destinationPath string, debugTemplates bool, delimLeft string, delimRight string) error {
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
	result, err := generateTemplate(s, filepath.Base(templatePath), delimLeft, delimRight)
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
