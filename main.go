package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

var branchName = os.Getenv("BRANCH_NAME")
var changeTarget = os.Getenv("CHANGE_TARGET")
var buildNumber = os.Getenv("BUILD_NUMBER")

func (v *Version) NextAlpha() (*Version, error) {
	var err error
	if v.packageName != "a" {
		return v, errors.New("package is not alpha")
	}
	v, err = v.calculateCommit()
	return v, err
}
func (v *Version) NextPatch() (*Version, error) {
	v.patch++
	return v, nil
}
func (v *Version) NextMinor() (*Version, error) {
	v.minor++
	return v, nil
}
func (v *Version) NextMajor() (*Version, error) {
	v.major++
	return v, nil
}

// Get current tags, find highest tag.
// get branch info.
// create version for every combination
func main() {
	tags := currentTags()
	version := getHighVersion(tags)
	nextVersion, err := version.guessNextVersion()
	if err != nil {
		panic(err)
	}

	gitUrl,_,_ := RunCommand("git","remote","get-url","origin")
	packageName := strings.Split(strings.Split(gitUrl,"/")[1],".git")[0]
	if strings.HasPrefix(gitUrl,"http") {
		packageName = strings.Split(strings.Split(gitUrl,"/")[4],".git")[0]
	}

	fmt.Println("Project Name:",packageName)
	fmt.Println("Git url:",gitUrl)
	fmt.Println("Running in",branchName, "branch with",version.Str(),"version")
	fmt.Println("Next version:",nextVersion.Str())
}
