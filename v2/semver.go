package main

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

// version must be formatted <major>.<minor>.<patch><separator><packageName><buildnumber>
// separator should be one of "+" or "-"
// buildnumber can be used total commit count of between last commit and last tagged version.
func CreateVersion(v string) (*Version, error) {
	var version = Version{}
	version.fullVersion = v
	if version.fullVersion == "" {
		return nil, errors.New("tag must not be empty")
	}

	// checks if string is correct
	const semanticRegex = `^(?P<major>0|[1-9]\d*)\.(?P<minor>0|[1-9]\d*)\.(?P<patch>0|[1-9]\d*)(?:-(?P<prerelease>(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+(?P<buildmetadata>[0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`
	r, _ := regexp.Compile(semanticRegex)
	if !r.MatchString(version.fullVersion) {
		return nil, errors.New("string is not a semantic version")
	}

	// splitting with dots, creates major minor patch
	splittedVersion := strings.Split(version.fullVersion, ".")
	if len(splittedVersion) < 2 {
		return nil, errors.New("tag must not be one character")
	}

	// checks if version has two chars
	// case will be implemented
	if len(splittedVersion) == 2 {
		version.isSemantic = false
	} else {
		version.isSemantic = true
	}

	if len(splittedVersion) > 3 {
		return nil, errors.New("tag has too many version. version must have two dots")
	}

	var err error
	version.major, err = strconv.Atoi(splittedVersion[0])
	if err != nil {
		return nil, errors.New("major version must be a number")
	}
	if version.major == 0 {
		version.state = "Initial Development"
	} else {
		version.state = "Production"
	}

	version.minor, err = strconv.Atoi(splittedVersion[1])
	if err != nil {
		return nil, errors.New("minor version must be a number")
	}

	afterMinor := strings.Join(splittedVersion[2:], "")

	// findNumStop finds where patch number finished.
	// You must separate your extra fields in patch with "+" or "-".
	findNumStop := func(str1 string) int {
		for key, char := range str1 {
			if _, err := strconv.Atoi(string(char)); err != nil {
				return key
			}
		}
		return 0
	}

	stopNo := findNumStop(afterMinor)

	if stopNo == 0 {
		// patch has only numeric chars.
		version.patch, err = strconv.Atoi(afterMinor)
		version.preRelease = 0
	} else {
		version.patch, err = strconv.Atoi(afterMinor[:stopNo])
		version.identifier = string(afterMinor[stopNo])
		// find package name
		re := regexp.MustCompile("[0-9]+")
		re2 := regexp.MustCompile("[A-Z|a-z]+")
		version.preRelease, _ = strconv.Atoi(re.FindAllString(afterMinor[stopNo+1:], -1)[0])
		version.packageName = re2.FindAllString(afterMinor[stopNo+1:], -1)[0]
	}

	return &version, nil
}

// Compares two version, returns highest version
func compare(v1, v2 string) (*Version, error) {
	version1, err := CreateVersion(v1)
	if err != nil {
		return nil, err
	}
	version2, err := CreateVersion(v2)
	if err != nil {
		return nil, err
	}

	compareInts := func(str1, str2 int) int {
		if str1 == str2 {
			return 0
		}
		if str1 > str2 {
			return 1
		}
		if str1 < str2 {
			return -1
		}
		return 0
	}
	returnVersion := func(int2 ...int) *Version {
		for i := range int2 {
			switch int2[i] {
			case 1:
				return version1
			case -1:
				return version2
			case 0:
				break
			}
		}
		return nil
	}

	resultMajor := compareInts(version1.major, version2.major)
	resultMinor := compareInts(version1.minor, version2.minor)
	resultPatch := compareInts(version1.patch, version2.patch)
	resultBuild := compareInts(version1.preRelease, version2.preRelease)

	if resultMajor == 0 && resultMinor == 0 && resultPatch == 0 {
		packageWeight := func(string2 string) int {
			if string2 == "a" || string2 == "alpha" {
				return 1
			}
			if string2 == "b" || string2 == "beta" {
				return 2
			}
			if string2 == "rc" {
				return 3
			}
			if string2 == "" {
				return 4
			}
			return 0
		}
		p1 := packageWeight(version1.packageName)
		p2 := packageWeight(version2.packageName)
		if p1 > p2 {
			return version1, nil
		} else if p2 >= p1 {
			return version2, nil
		}
	}

	return returnVersion(resultMajor, resultMinor, resultPatch, resultBuild), nil
}
