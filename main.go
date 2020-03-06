package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"syscall"
)

type (
	Version struct {
		major int
		minor int
		patch int
		preRelease int
		packageName string
		err        error
	}
)


var BRANCH_NAME = os.Getenv("BRANCH_NAME")
var CHANGE_TARGET = os.Getenv("CHANGE_TARGET")

func main()  {

	tags := getTags()
	HighTag := getHighVersion(tags)
	vNext := guessNextVersion(HighTag)
	vNext = vNext.NextVersion()
	fmt.Println("Repository version: ",rtnVersion(HighTag))
	fmt.Println("Suggested next patched version: ",rtnVersion(HighTag.NextVersion()))

	createEnvs(rtnVersion(vNext))

}

func errCheck(err error,errorString string) {
	if err != nil{
		log.Fatal(errorString)
	}
}

func getHighVersion(tags []string) Version {
	var HighTag = Version{}
	for i := 0; i < len(tags); i++ {
		v1, ok := ParseVersion(tags[i])
		if ok == false {
			continue
		}
		if i == 0 {
			HighTag = v1
			continue
		}
		v2, ok := ParseVersion(tags[i-1])
		if ok == false {
			continue
		}
		HighTag = Compare(v1,v2)
	}
	return HighTag
}

func getTags() []string {
	dir, err := os.Getwd()
	errCheck(err,"Dizin alınamadı")

	var gitTags []string
	validID := regexp.MustCompile(`.*refs/tags/.*`)
	outputCommand, _,_ := RunCommand("git","--git-dir",dir + "/.git", "show-ref", "--tags")

	ifMatched := validID.MatchString(outputCommand)
	if ifMatched {
		parseStrings := strings.Split(outputCommand,"\n")
		for _,v := range parseStrings {
			if v == "" {
				continue
			}
			parseStringsNew := strings.Split(v,"refs/tags/")
			gitTags = append(gitTags,parseStringsNew[1])
		}
	} else {
		// yoksa ilk versiyon veriliyor.
		gitTags = append(gitTags,"0.1.0")
		return gitTags
	}
	return gitTags
}

func rtnVersion(v Version) string {

	a := strconv.Itoa(v.major)
	b := strconv.Itoa(v.minor)
	c := strconv.Itoa(v.patch)
	d := v.packageName
	e := strconv.Itoa(v.preRelease)

	if v.packageName == "" {
		return a + "." + b + "." + c
	}
	return a + "." + b + "." + c  + d + e
}

func RunCommand(name string, args ...string) (stdout string, stderr string, exitCode int) {
	var outbuf, errbuf bytes.Buffer
	cmd := exec.Command(name, args...)
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf
	err := cmd.Run()
	stdout = outbuf.String()
	stderr = errbuf.String()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			ws := exitError.Sys().(syscall.WaitStatus)
			exitCode = ws.ExitStatus()
		} else {
			exitCode = 1
			if stderr == "" {
				stderr = err.Error()
			}
		}
	} else {
		ws := cmd.ProcessState.Sys().(syscall.WaitStatus)
		exitCode = ws.ExitStatus()
	}
	return
}

func createEnvs(v string)  {
	gitUrl,_,_ := RunCommand("git","remote","get-url","origin")
	packageName := strings.Split(strings.Split(gitUrl,"/")[1],".git")[0]
	fmt.Println("PACKAGE_NAME="+packageName)
	fmt.Println("PACKAGE_VERSION="+v)
	fmt.Println("GIT_URL="+gitUrl)
	f, err := os.Create("BUILD_CONTEXT_FILE")
	if err != nil {
		fmt.Println(err)
	}
	writeFile := []string{
		"PACKAGE_NAME: " + packageName,
		"PACKAGE_VERSION: " + v,
		"GIT_URL: " + gitUrl,
	}
	for _, v := range writeFile {
		_, _ = fmt.Fprintln(f, v)
	}
	_ = f.Close()
}
