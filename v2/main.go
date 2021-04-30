package main

import (
	"errors"
	"fmt"
	"os"
)

var branchName = os.Getenv("BRANCH_NAME")
var changeTarget = os.Getenv("CHANGE_TARGET")
var buildNumber = os.Getenv("BUILD_NUMBER")




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



func main() {
	version, err := CreateVersion("1.1.1-b12")
	if err != nil {
		panic(err)
	}
	asd, _ := version.calculateCommit()
	fmt.Println(version.Str(),asd)
}
