package main

import (
	"github.com/flosch/pongo2"
	"os"
	"strings"
)

type PongoTemplar struct {
	source string
}

func (templar *PongoTemplar) generateTemplate() (string, error) {
	context := pongo2.Context{}

	tmpl, err := pongo2.FromString(templar.source)
	if err != nil {
		return "", err
	}

	for _, val := range os.Environ() {
		parts := strings.SplitN(val, "=", 2)
		key, value := parts[0], parts[1]

		context[key] = value
	}

	out, err := tmpl.Execute(context)
	if err != nil {
		return "", err
	}

	return out, nil
}
