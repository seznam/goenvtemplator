package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

type Templar interface {
	generateTemplate() (string, error)
}

type templatePaths struct {
	source      string
	destination string
}

func (t templatePaths) String() string {
	return fmt.Sprintf("{source: '%s', destination: '%s'}",
		t.source, t.destination)
}

type templatesPaths []templatePaths

func (ts *templatesPaths) Set(value string) error {
	var t templatePaths
	parts := strings.Split(value, ":")
	if len(parts) == 2 {
		t.source = strings.TrimSpace(parts[0])
		t.destination = strings.TrimSpace(parts[1])
	} else {
		return errors.New("Option has invalid format!")
	}
	*ts = append(*ts, t)
	return nil
}

func (ts *templatesPaths) String() string {
	return fmt.Sprintf("%v", *ts)
}

// to parse slice of strings from flags we need to use custom type
type envFiles []string

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

func generateTemplates(
	ts templatesPaths, debug bool, delimLeft string, delimRight string, engine string) error {
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
		default:
			templar = &TextTemplar{
				source:     source,
				name:       templateName,
				delimLeft:  delimLeft,
				delimRight: delimRight,
			}
		}

		render, err := templar.generateTemplate()
		if err != nil {
			return err
		}

		if err = writeFile(t.destination, render); err != nil {
			return err
		}

	}
	return nil
}

var (
	v            int
	buildVersion string = "Build version was not specified."
)

func main() {
	var tmpls templatesPaths
	flag.Var(&tmpls, "template", "Template (/template:/dest). Can be passed multiple times.")
	var debugTemplates bool
	flag.BoolVar(&debugTemplates, "debug-templates", false, "Print processed templates to stdout.")
	var doExec bool
	flag.BoolVar(&doExec, "exec", false, "Activates exec by command. First non-flag arguments is the command, the rest are it's arguments.")
	var printVersion bool
	flag.BoolVar(&printVersion, "version", false, "Prints version.")
	var envFileList envFiles
	flag.Var(&envFileList, "env-file", "Additional file with environment variables. Can be passed multiple times.")
	var delimLeft string
	flag.StringVar(&delimLeft, "delim-left", "", "Override default left delimiter {{.")
	var delimRight string
	flag.StringVar(&delimRight, "delim-right", "", "Override default right delimiter }}.")
	flag.IntVar(&v, "v", 0, "Verbosity level.")

	flag.Parse()

	// if no env-file was passed, godotenv.Load loads .env file by default, we want to disable this
	if len(envFileList) > 0 {
		err := godotenv.Load(envFileList...)
		if err != nil {
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

	if err := generateTemplates(tmpls, debugTemplates, delimLeft, delimRight, "default"); err != nil {
		panic(err)
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
