package main

import (
	"log"
	"strconv"
	"strings"
)

func (v Version) calculateCommit() Version {
	v1 := rtnVersion(v)
	//dir, _ := os.Getwd()
	commitCount, _, _ := RunCommand("git", "rev-list",v1 + "..HEAD" , "--count")
	c1, _ := strconv.Atoi(strings.Replace(commitCount,"\n","",-1))
	v.preRelease = c1
	return v
}

func guessNextVersion(v Version) Version  {
	v = v.calculateCommit()
	switch {
	case strings.HasPrefix(BRANCH_NAME,"feature"):
		/// Next feature version
		v.packageName = "a"
		return v
	case strings.HasPrefix(BRANCH_NAME,"release"):
		// Next Release Candidate version
		v.packageName = "rc"
		return v
	case  strings.HasPrefix(BRANCH_NAME,"PR"):
		//
		if CHANGE_TARGET == "master" {
			vRelease := findReleases()
			return vRelease
		}
		if CHANGE_TARGET == "develop" {
			v.packageName = "a"
			return v
		}
	case strings.HasPrefix(BRANCH_NAME,"develop"):
		v.packageName = "b"
		return v
	case strings.HasPrefix(BRANCH_NAME,"master"):
		return findReleases()
	}
	return v
}

func (v Version) NextVersion() Version {
	v.patch -= -1 //swh
	return v
}

func findReleases()  Version {
	ReleaseTags,_,_ := RunCommand("git","branch","-r")
	var releaseBranch []string
	var release Version
	var hotfixBranch []string
	var hotfix Version
	for _,i := range strings.Split(ReleaseTags,"\n") {
		i = strings.Join(strings.Split(i, "/")[1:], "/")
		if strings.HasPrefix(i, `release`) {
			releaseBranch = append(releaseBranch, strings.Join(strings.Split(i, "/")[1:], ""))
		}
		if strings.HasPrefix(i,`hotfix`) {
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
		return Version{}
	}


}