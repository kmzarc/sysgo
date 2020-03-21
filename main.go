package main

import (
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	execStart   = kingpin.Flag("exec", "command to run including parameters").Required().Short('e').String()
	workingDir  = kingpin.Flag("dir", "set the working directory for the process").Short('d').String()
	waitBetween = kingpin.Flag("wait", "How long to wait between crashes before try again").Default("1s").Short('w').Duration()
	count       = kingpin.Flag("count", "Number of retry if process fails. 0 means no limit").Default("2").Short('c').Int()
)

func main() {
	kingpin.Version("0.0.1")
	kingpin.Parse()

	log.Printf("Starting %s\n", *execStart)
	for i := 0; i < *count || *count == 0; i++ {
		runProcess()
		if i < *count || *count == 0 {
			log.Printf("waiting for for next run %s\n", *waitBetween)
			time.Sleep(*waitBetween)
		}
	}
}

func runProcess() {
	name, args := parseExecFlag(*execStart)
	cmd := exec.Command(name, args...)
	cmd.Dir = *workingDir

	// use os.Stdout and osStderr for the command
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		log.Printf("failed to start %s\n", err)
		return
	}
	log.Printf("Running process: %d\n", cmd.Process.Pid)

	err = cmd.Wait()
	if err != nil {
		log.Printf("command existed with error %s\n", err)
		return
	}

	log.Printf("command existed\n")
}

// parse the exec command to what exec.Cmd needs
func parseExecFlag(s string) (string, []string) {
	l := strings.FieldsFunc(s, splitFunc)
	return l[0], l[1:]
}

var last rune

// split need to keep the state of parsing the commands
// this function is not a complete parses for commands
func splitFunc(r rune) bool {
	if last == '"' && r == '"' {
		last = '0'
		return r == '"'
	}
	if last == '"' && r != '"' {
		return r == '"'
	}
	last = r
	return r == ' ' || r == '"'
}
