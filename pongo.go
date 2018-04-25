package main

import (
	"os"
	"strings"

	"github.com/flosch/pongo2"
)

type PongoTemplar struct {
	Source string
}

func (templar *PongoTemplar) generateTemplate() (string, error) {
	context := pongo2.Context{}

	tmpl, err := pongo2.FromString(templar.Source)
	if err != nil {
		return "", err
	}

	for _, val := range os.Environ() {
		parts := strings.SplitN(val, "=", 2)
		key, value := parts[0], parts[1]

		context[key] = value
	}

	Debug("Using context %v", context)

	out, err := tmpl.Execute(context)
	if err != nil {
		return "", err
	}

	return out, nil
}
