package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/joho/godotenv"
)

var (
	v            int
	buildVersion string = "Build version was not specified."
	DEBUG        bool
)

type templatesPaths []templatePath

// to parse slice of strings from flags we need to use custom type
type envFiles []string

type Templar interface {
	generateTemplate() (string, error)
}

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
		log.Fatalf("Destination path '%s' is not absolute!", destinationPath)
		return errors.New("absolute path error")
	}

	if err := ioutil.WriteFile(destinationPath, []byte(data), 0664); err != nil {
		return err
	}

	return nil
}

func readSource(templatePath string) (string, error) {
	if !filepath.IsAbs(templatePath) {
		log.Fatalf("Template path '%s' is not absolute!", templatePath)
		return "", errors.New("absolute path error")
	}

	var slice []byte
	var err error
	if slice, err = ioutil.ReadFile(templatePath); err != nil {
		return "", err
	}

	return string(slice), nil

}

func Debug(message string, args ...interface{}) {
	if DEBUG {
		log.Printf(message, args...)
	}
}

func generateTemplates(
	ts templatesPaths,
	delimLeft string,
	delimRight string,
	engine string) error {

	for _, t := range ts {
		if v > 0 {
			log.Printf("generating %s -> %s", t.source, t.destination)
		}

		var templar Templar

		source, err := readSource(t.source)
		if err != nil {
			return err
		}

		templateName := filepath.Base(t.source)

		switch engine {
		case "pongo":
			templar = &PongoTemplar{
				Source: source,
			}
		case "text/template":
			templar = &TextTemplar{
				Source:     source,
				Name:       templateName,
				DelimLeft:  delimLeft,
				DelimRight: delimRight,
			}
		}

		Debug("Templating %s", templateName)

		render, err := templar.generateTemplate()
		if err != nil {
			return err
		}

		Debug("Generated template %s", render)

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
	var delimLeft string
	var delimRight string
	var engine string

	flag.Var(&tmpls, "template", "Template (/template:/dest). Can be passed multiple times.")
	flag.BoolVar(&DEBUG, "debug-templates", false, "Print processed templates to stdout.")
	flag.BoolVar(&doExec, "exec", false, "Activates exec by command. First non-flag arguments is the command, the rest are it's arguments.")
	flag.BoolVar(&printVersion, "version", false, "Prints version.")
	flag.Var(&envFileList, "env-file", "Additional file with environment variables. Can be passed multiple times.")
	flag.StringVar(&delimLeft, "delim-left", "", "Override default left delimiter {{.")
	flag.StringVar(&delimRight, "delim-right", "", "Override default right delimiter }}.")
	flag.IntVar(&v, "v", 0, "Verbosity level.")
	flag.StringVar(&engine, "engine", "text/template", "Override default text/template [support: pongo2]")

	flag.Parse()
	// if no env-file was passed, godotenv.Load loads .env file by default, we want to disable this
	if len(envFileList) > 0 {
		if err := godotenv.Load(envFileList...); err != nil {
			log.Fatal(err)
		}
	}

	if printVersion {
		log.Printf("Version: %s", buildVersion)
		os.Exit(0)
	}

	if v > 0 {
		log.Print("Generating templates")
	}

	if err := generateTemplates(tmpls, delimLeft, delimRight, engine); err != nil {
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
