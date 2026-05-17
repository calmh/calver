package main

import (
	"flag"
	"fmt"
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
	curYear := now.Year() - 2000
	curMonth := now.Month()

	if prev == "" {
		prev, _ = cmd("git", "describe", "--abbrev=0", "--match", "v[0-9].*")
	}
	prev = strings.TrimLeft(prev, "vV")

	if prev == "" { // no previous version
		fmt.Printf("%d.%d.0\n", curYear, curMonth)
		return
	}

	latestStable, err := semver.NewVersion(prev)
	if err != nil { // unparseable version
		fmt.Printf("%d.%d.0\n", curYear, curMonth)
		return
	}

	if latestStable.Major != int64(curYear) || latestStable.Minor != int64(curMonth) { // new month/year
		fmt.Printf("%d.%d.0\n", curYear, curMonth)
		return
	}

	latestStable.Patch++
	fmt.Println(latestStable)
}

func cmd(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	bs, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(bs)), nil
}
