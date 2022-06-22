package build

import (
	"log"
	"os/exec"
	"strings"
	"time"
)

var mode = "development"

var commit = "<unknown>"
var version = "<unknown>"
var branch = "<unknown>"
var buildTime = "<unknown>"

var goTime time.Time

func init() {
	if mode == "development" {
		commit = "dev"
		version = "dev"

		cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
		stdout, err := cmd.Output()
		if err != nil {
			log.Println("failed to get git branch:", err)
		} else {
			branch = strings.Replace(string(stdout), "\n", "", -1)
		}

		goTime = time.Now()
	}
	if buildTime == "<unknown>" {
		goTime = time.Now()
	} else {
		var err error
		goTime, err = time.Parse(time.RFC3339, buildTime)
		if err != nil {
			log.Printf("[error] failed to parse build time: %v\n", err)
			goTime = time.Now()
		}
	}
}

func Mode() string    { return mode }
func Commit() string  { return commit }
func Version() string { return version }
func Branch() string  { return branch }
func Time() time.Time { return goTime }
