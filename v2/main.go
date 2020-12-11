package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var branchName = os.Getenv("BRANCH_NAME")
var changeTarget = os.Getenv("CHANGE_TARGET")

type Version struct {
	major       int
	minor       int
	patch       int
	preRelease  int
	packageName string
	err         error
	isSemantic	bool
	fullVersion string
}

func CreateVersion(v string) (*Version,error) {
	var version = Version{}
	version.fullVersion = v
	if v == "" {
		return nil,errors.New("tag must not be empty")
	}

	// splitting with dots, creates major minor patch
	splittedV := strings.Split(v,".")
	if len(splittedV) < 2 {
		return nil,errors.New("tag must not be one character")
	}

	// checks if version has two chars
	// case will be implemented
	if len(splittedV) == 2 {
		version.isSemantic = false
	} else {
		version.isSemantic = true
	}

	if len(splittedV) > 3 {
		return nil,errors.New("tag has too many version. version should have two dots")
	}


	version.major, _ = strconv.Atoi(splittedV[0])
	version.minor, _ = strconv.Atoi(splittedV[1])

	// checks if patch has package classifier
	const alpha = "abcdefghijklmnopqrstuvwxyz"
	isPackage := func() bool {
		for _,char := range splittedV[2] {
			if strings.ContainsAny(alpha,strings.ToLower(string(char))) {
				return true
			}
		}
		return false
	}
	if isPackage() {
		if strings.Index(splittedV[2],"b") != -1 {
			version.packageName = "b"
		} else if strings.Index(splittedV[2],"a") != -1 {
			version.packageName = "a"
		} else if strings.Index(splittedV[2],"rc") != -1 {
			version.packageName = "rc"
		}
		version.patch, _ = strconv.Atoi(strings.Split(splittedV[2], version.packageName)[0])
		version.preRelease, _ = strconv.Atoi(strings.Split(splittedV[2], version.packageName)[1])
	} else {
		version.patch, _ = strconv.Atoi(splittedV[2])
		version.preRelease = 0
	}
	return &version,nil
}

func (v Version) NextAlpha() error {
	if v.packageName != "a" {
		return errors.New("package is not alpha")
	}
	//calculate commits

	return nil
}
func (v Version) NextPatch() error {

	return nil
}
func (v Version) NextMinor() error {

	return nil
}
func (v Version) NextMajor() error {

	return nil
}


func (v Version) current()  {
	dir, err := os.Getwd()
	if err != nil {
		log.Panic("dizin alinamadi" , err)
	}
	validID := regexp.MustCompile(`.*refs/tags/.*`)
	outputCommand, _, _ := RunCommand("git", "--git-dir", dir+"/.git", "show-ref", "--tags")
	ifMatched := validID.MatchString(outputCommand)
	var gitTags = []string{}
	if ifMatched {
		parseStrings := strings.Split(outputCommand, "\n")
		for _, v := range parseStrings {
			if v == "" {
				continue
			}
			parseStringsNew := strings.Split(v, "refs/tags/")
			gitTags = append(gitTags, parseStringsNew[1])
		}
	} else {
		gitTags = append(gitTags, "0.1.0")
	}
}

func (v Version) parse(str string) error {
	var splitWithDot = strings.Split(str, ".")
	v.major, _ = strconv.Atoi(splitWithDot[0])
	v.minor, _ = strconv.Atoi(splitWithDot[1])
	return nil
}

func main() {
	version, _ := CreateVersion("1.1.1b123")
	fmt.Println(version)
}