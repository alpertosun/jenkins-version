package main

import (
	"log"
	"strconv"
	"strings"
)

func boolCheck(bool2 bool,string2 string)  {
	if !bool2 {
		log.Fatal(string2)
		return
	}
}


func CompareTwoVersion(first,second string) Version {
	version1, _ := ParseVersion(first)
	version2, _ := ParseVersion(second)
	version := Compare(version1,version2)
	return version
}

func Compare(version1 Version, version2 Version) Version {

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
	
	returnVersion := func(int2 ...int) (Version) {
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
		return Version{}
	}

	resultMajor := compareInts(version1.major, version2.major)
	resultMinor := compareInts(version1.minor, version2.minor)
	resultPatch := compareInts(version1.patch, version2.patch)
	resultBuild := compareInts(version1.preRelease,version2.preRelease)


	if resultMajor == 0 && resultMinor == 0 && resultPatch == 0 {
		packageWeight := func(string2 string) int {
			if string2 == "a" {
				return 1
			}
			if string2 == "b" {
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
			return version1
		} else if p2 >= p1 {
			return version2
		}
	}

	return returnVersion(resultMajor,resultMinor,resultPatch,resultBuild)
}



func ParseVersion(tag string) (Version,bool) {
	version := Version{}
	if tag == "" {
		log.Fatal("Tag boş olamaz.")
		return version,false
	}

	var splitWithDot = strings.Split(tag, ".")
	version.major, version.err = strconv.Atoi(splitWithDot[0])
	if version.err != nil {
		return Version{},false
	}
	version.minor, version.err = strconv.Atoi(splitWithDot[1])
	if version.err != nil {
		return Version{},false
	}
	if strings.Index(splitWithDot[2],"b") != -1 {
		version.packageName = "b"
	} else if strings.Index(splitWithDot[2],"a") != -1 {
		version.packageName = "a"
	} else if strings.Index(splitWithDot[2],"rc") != -1 {
		version.packageName = "rc"
	}

	if version.packageName == "" {
		version.patch, version.err = strconv.Atoi(splitWithDot[2])
		errCheck(version.err,"Yanlış paket ismi kullanılmış.")
		return version,true
	} else {
		packageBeta := strings.Split(splitWithDot[2],version.packageName)
		version.patch, _ = strconv.Atoi(packageBeta[0])
		version.preRelease, _ = strconv.Atoi(packageBeta[1])
	}
	return version,true
}

