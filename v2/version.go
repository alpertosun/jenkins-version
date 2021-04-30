package main

import (
	"regexp"
	"strconv"
	"strings"
)

type Version struct {
	major       int
	minor       int
	patch       int
	preRelease  int
	packageName string
	err         error
	isSemantic  bool
	state       string
	fullVersion string
	identifier  string
}

func (v Version) calculateCommit() (*Version, error) {
	var err error
	currentVersion := ""
	commitCount, _, _ := RunCommand("git", "rev-list", currentVersion+"..HEAD", "--count")
	v.preRelease, err = strconv.Atoi(strings.Replace(commitCount, "\n", "", -1))
	return &v, err
}

func (v Version) currentTags() []string {
	validID := regexp.MustCompile(`.*refs/tags/.*`)
	outputCommand, _, _ := RunCommand("git", "show-ref", "--tags")
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
	return gitTags
}

func (v *Version) Str() string {
	return strconv.Itoa(v.major) + "." + strconv.Itoa(v.minor) + "." + strconv.Itoa(v.patch) + v.identifier + v.packageName + strconv.Itoa(v.preRelease)
}
