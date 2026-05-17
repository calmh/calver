package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/coreos/go-semver/semver"
)

func main() {
	var prev string
	flag.StringVar(&prev, "prev", "", "Previous version, leave blank to detect")
	flag.Parse()

	now := time.Now().UTC()

	ver := version(prev, now)
	fmt.Println(ver)

	if out := os.Getenv("GITHUB_OUTPUT"); out != "" {
		fd, err := os.OpenFile(out, os.O_WRONLY|os.O_APPEND, 0o644)
		if err != nil {
			fmt.Printf("Failed to open $GITHUB_OUTPUT: %v\n", err)
			os.Exit(1)
		}
		_, _ = fmt.Fprintf(fd, "version=%s\n", ver)
		_ = fd.Close()
	}
}

func version(prev string, now time.Time) string {
	if prev == "" {
		var err error
		prev, err = cmd("git", "describe", "--abbrev=0", "--tags", "--match", "v[0-9]*")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to run git: %v\n", err)
			fmt.Fprintf(os.Stderr, "output: %v\n", prev)
			prev = ""
		}
	}
	prev = strings.TrimLeft(prev, "vV")

	curYear := now.Year() - 2000
	curMonth := now.Month()

	if prev == "" { // no previous version
		return fmt.Sprintf("%d.%d.0", curYear, curMonth)
	}

	latestStable, err := semver.NewVersion(prev)
	if err != nil { // unparseable version
		return fmt.Sprintf("%d.%d.0", curYear, curMonth)
	}

	if latestStable.Major != int64(curYear) || latestStable.Minor != int64(curMonth) { // new month/year
		return fmt.Sprintf("%d.%d.0", curYear, curMonth)
	}

	latestStable.Patch++
	return latestStable.String()
}

func cmd(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	bs, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(bs)), err
}
