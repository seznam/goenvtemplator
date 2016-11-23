package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

type templatePaths struct {
	source      string
	destination string
}

func (t templatePaths) String() string {
	return fmt.Sprintf("{source: '%s', destination: '%s'}", t.source, t.destination)
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

func generateTemplates(ts templatesPaths, debug bool) error {
	for _, t := range ts {
		if v > 0 {
			log.Printf("generating %s -> %s", t.source, t.destination)
		}
		if err := generateFile(t.source, t.destination, debug); err != nil {
			return fmt.Errorf("Error while generating '%s' -> '%s'. %v", t.source, t.destination, err)
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

	if err := generateTemplates(tmpls, debugTemplates); err != nil {
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
