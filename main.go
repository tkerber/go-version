package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// cmd executes a command and returns its trimmed output, or the empty string
// if an error occurred.
func cmd(str string) string {
	split := strings.Split(str, " ")
	ret, err := exec.Command(split[0], split[1:]...).Output()
	if err != nil {
		return ""
	} else {
		return strings.Trim(string(ret), " \t\n\r")
	}
}

// generateGit generates the version file for git.
func generateGit(pkg string, out io.Writer) {
	shortHash := cmd("git rev-parse --short HEAD")
	longHash := cmd("git rev-parse HEAD")
	tag := cmd("git describe --abbrev=0 --tags")
	dateStr := cmd("git show -s --format=%ct")
	dateUnix, err := strconv.ParseInt(dateStr, 0, 64)
	if err != nil {
		dateUnix = 0
	}
	out.Write([]byte(fmt.Sprintf(`package %s

import "time"

const (
	CommitHashShort = %#v
	CommitHashLong  = %#v
	CommitTag       = %#v
)

var CommitDate = time.Unix(%#v, 0)
`, pkg, shortHash, longHash, tag, dateUnix)))
}

// generateBzr generates the version file for bazaar.
func generateBzr(pkg string, out io.Writer) {
	dateStr := cmd("bzr version-info --custom --template {date}")
	date, err := time.Parse("2006-01-02 15:04:05 -0700", dateStr)
	if err != nil {
		date = time.Unix(0, 0)
	}
	revnoStr := cmd("bzr version-info --custom --template {revno}")
	revno, err := strconv.ParseInt(revnoStr, 0, 64)
	if err != nil {
		revno = 0
	}
	revisionId := cmd("bzr version-info --custom --template {revision_id}")
	tagsStr := cmd("bzr tags --sort=time")
	tags := strings.SplitN(tagsStr, "\n", 2)
	// First group is the tag name
	re := regexp.MustCompile(`^(.*)\s+\S+$`)
	match := re.FindStringSubmatch(tags[0])
	tag := ""
	if match != nil && len(match) >= 2 {
		tag = match[1]
	}
	out.Write([]byte(fmt.Sprintf(`package %s

import "time"

const (
	RevNo      = %#v
	RevisionId = %#v
	CommitTag  = %#v
)

var CommitDate = time.Unix(%#v, 0)
`, pkg, revno, revisionId, tag, date.Unix())))
}

// generateHg generates the version file for mercurial.
func generateHg(pkg string, out io.Writer) {
	dateStr := cmd("hg heads . -T {date}")
	dateUnix, err := strconv.ParseFloat(dateStr, 64)
	if err != nil {
		dateUnix = 0
	}
	tag := cmd("hg heads . -T {latesttag}")
	revStr := cmd("hg heads . -T {rev}")
	revNo, err := strconv.ParseInt(revStr, 0, 64)
	if err != nil {
		revNo = 0
	}
	hash := cmd("hg heads . -T {node}")
	out.Write([]byte(fmt.Sprintf(`package %s

import "time"

const (
	CommitTag  = %#v
	RevNo      = %#v
	CommitHash = %#v
)

var CommitDate = time.Unix(%#v, 0)
`, pkg, tag, revNo, hash, int(dateUnix))))
}

// A repo is a type of repository.
type repo int

// The supported types of repositories.
const (
	git repo = iota
	bzr
	hg
)

func repoType(dir string) (repo, bool) {
	if _, err := os.Stat(filepath.Join(dir, ".git")); !os.IsNotExist(err) {
		return git, true
	} else if _, err := os.Stat(filepath.Join(dir, ".bzr")); !os.IsNotExist(err) {
		return bzr, true
	} else if _, err := os.Stat(filepath.Join(dir, ".bzr")); !os.IsNotExist(err) {
		return hg, true
	} else {
		newDir := filepath.Dir(dir)
		if newDir != dir {
			return repoType(newDir)
		} else {
			return 0, false
		}
	}
}

func main() {
	outfile := flag.String("o", "./version.go", "The output file to generate.")
	pkg := flag.String("pkg", "main", "The package of the output file.")
	flag.Parse()
	wd, err := os.Getwd()
	if err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}
	f, err := os.Create(*outfile)
	if err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}
	defer f.Close()
	repo, ok := repoType(wd)
	if !ok {
		fmt.Fprintf(os.Stderr, "No repository found.\n")
		f.Close()
		os.Exit(1)
	}
	switch repo {
	case git:
		generateGit(*pkg, f)
	case bzr:
		generateBzr(*pkg, f)
	case hg:
		generateHg(*pkg, f)
	}
}
