package main

import (
	"log"
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
	currentVersion := v.Str()
	commitCount, _, _ := RunCommand("git", "rev-list", currentVersion+"..HEAD", "--count")
	v.preRelease, err = strconv.Atoi(strings.Replace(commitCount, "\n", "", -1))
	v.preRelease++
	return &v, err
}

func currentTags() []string {
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

func (v Version) currentReleases() *Version {
	releaseTags, _, _ := RunCommand("git", "branch", "-r")
	var releaseBranch []string
	var release *Version
	var hotfixBranch []string
	var hotfix *Version
	for _, i := range strings.Split(releaseTags, "\n") {
		i = strings.Join(strings.Split(i, "/")[1:], "/")
		if strings.HasPrefix(i, `release`) {
			releaseBranch = append(releaseBranch, strings.Join(strings.Split(i, "/")[1:], ""))
		}
		if strings.HasPrefix(i, `hotfix`) {
			hotfixBranch = append(hotfixBranch, strings.Join(strings.Split(i, "/")[1:], ""))
		}
	}
	if len(releaseBranch) == 0 && len(hotfixBranch) != 0 {
		hotfix = getHighVersion(hotfixBranch)
		return hotfix

	} else if len(releaseBranch) != 0 && len(hotfixBranch) == 0 {
		release = getHighVersion(releaseBranch)
		return release
	} else {
		log.Fatal("No release & hotfix branches found")
		return nil
	}
}

func (v *Version) Str() string {
	if v.preRelease == 0 {
		return strconv.Itoa(v.major) + "." + strconv.Itoa(v.minor) + "." + strconv.Itoa(v.patch)
	}
	return strconv.Itoa(v.major) + "." + strconv.Itoa(v.minor) + "." + strconv.Itoa(v.patch) + "-" + v.identifier + v.packageName + strconv.Itoa(v.preRelease)
}

func (v *Version) guessNextVersion() (*Version, error) {
	v, err := v.calculateCommit()
	branchInfo := func(branchPrefix string) bool {
		return strings.HasPrefix(branchName, branchPrefix)
	}
	switch {
	case branchInfo("feature"):
		v.packageName = "alpha"
	case branchInfo("PR"):
		if changeTarget == "master" {
			//vRelease := findReleases()
			//return vRelease
		}
		if changeTarget == "develop" {
			v.packageName = "alpha"
		}
	case branchInfo("develop"):
		v.packageName = "beta"
	case branchInfo("master"):
		v.packageName = ""
		v.preRelease = 0
	case branchInfo("release"):
		if changeTarget == "master" {
			v.packageName = "rc"
		}
		if changeTarget == "develop" {
			v.packageName = "beta"
		}
	case branchInfo("hotfix"):
		if changeTarget == "master" {
			v.packageName = ""
			v, err = v.NextPatch()
		}
		if changeTarget == "develop" {
			v.packageName = "beta"
			v, err = v.NextPatch()
		}
	}
	v, err = v.NextPatch()
	return v, err
}
