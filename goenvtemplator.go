package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	log "github.com/golang/glog"
	"github.com/joho/godotenv"

	"github.com/seznam/goenvtemplator/engine"
)

var (
	buildVersion string = "Build version was not specified."
)

type templatesPaths []templatePath

// to parse slice of strings from flags we need to use custom type
type envFiles []string

type templatePath struct {
	source      string
	destination string
}

func (t templatePath) String() string {
	return fmt.Sprintf("{source: '%s', destination: '%s'}",
		t.source, t.destination)
}

func (ts *templatesPaths) Set(value string) error {
	parts := strings.SplitN(value, ":", 2)
	if len(parts) < 2 {
		return errors.New("Option has invalid format!")
	}

	*ts = append(*ts, templatePath{
		source:      strings.TrimSpace(parts[0]),
		destination: strings.TrimSpace(parts[1]),
	})

	return nil
}

func (ts *templatesPaths) String() string {
	return fmt.Sprintf("%v", *ts)
}

func (ef *envFiles) Set(value string) error {
	*ef = append(*ef, value)
	return nil
}

func (ef *envFiles) String() string {
	return fmt.Sprintf("%v", *ef)
}

func writeFile(destinationPath string, data string) error {
	if !filepath.IsAbs(destinationPath) {
		return fmt.Errorf("absolute path error: %s", destinationPath)
	}

	if err := ioutil.WriteFile(destinationPath, []byte(data), 0664); err != nil {
		return err
	}

	return nil
}

func readFile(templatePath string) (string, error) {
	if !filepath.IsAbs(templatePath) {
		return "", fmt.Errorf("absolute path error: %s", templatePath)
	}

	var slice []byte
	var err error
	if slice, err = ioutil.ReadFile(templatePath); err != nil {
		return "", err
	}

	return string(slice), nil

}

func generateTemplates(ts templatesPaths, engineName string) error {
	for _, t := range ts {
		if log.V(1) {
			log.Info("generating %s -> %s", t.source, t.destination)
		}

		var templar engine.Templar

		source, err := readFile(t.source)
		if err != nil {
			return err
		}

		templateName := filepath.Base(t.source)

		switch engineName {
		case "pongo":
			templar = &engine.PongoTemplar{
				Source: source,
			}
		case "text/template":
			templar = &engine.TextTemplar{
				Source: source,
				Name:   templateName,
			}
		}

		if log.V(3) {
			log.Info("Templating %s", templateName)
		}

		render, err := templar.GenerateTemplate()
		if err != nil {
			return err
		}

		if log.V(3) {
			log.Info("Generated template %s", render)
		}

		if err = writeFile(t.destination, render); err != nil {
			return err
		}

	}
	return nil
}

func main() {
	var tmpls templatesPaths
	var doExec bool
	var printVersion bool
	var envFileList envFiles
	var engine string

	flag.Var(&tmpls, "template", "Template (/template:/dest). Can be passed multiple times.")
	flag.BoolVar(&doExec, "exec", false, "Activates exec by command. First non-flag arguments is the command, the rest are it's arguments.")
	flag.BoolVar(&printVersion, "version", false, "Prints version.")
	flag.Var(&envFileList, "env-file", "Additional file with environment variables. Can be passed multiple times.")
	flag.StringVar(
		&engine, "engine", "text/template",
		"Override default text/template [supports: text/template, pongo]",
	)

	flag.Parse()
	// if no env-file was passed, godotenv.Load loads .env file by default, we want to disable this
	if len(envFileList) > 0 {
		if err := godotenv.Load(envFileList...); err != nil {
			log.Fatal(err)
		}
	}

	if printVersion {
		log.Info("Version: %s", buildVersion)
		os.Exit(0)
	}

	if log.V(1) {
		log.Info("Generating templates")
	}

	if err := generateTemplates(tmpls, engine); err != nil {
		log.Fatal(err)
	}

	if doExec {
		if flag.NArg() < 1 {
			log.Fatal("Missing command to execute!")
		}

		args := flag.Args()
		cmd := args[0]
		cmdPath, err := exec.LookPath(cmd)
		if err != nil {
			log.Fatal(err)
		}
		if err := syscall.Exec(cmdPath, args, os.Environ()); err != nil {
			log.Fatalf("Unable to exec '%s'! %v", cmdPath, err)
		}
	}
}
