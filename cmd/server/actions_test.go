package main

import (
	"testing"
)

// To run this tests the directories on the currentPath have to be created, including
// the base directory

func TestActionMkdir(t *testing.T) {
	const currentPath = "/testdir/testdir2"
	err := mkdir(currentPath, "~/test", "~/.file", "helloworld/")
	if err != nil {
		t.Error(err)
	}
}

func TestActionCd1(t *testing.T) {
	const currentPath = "/testdir/testdir2"
	newPath, err := cd(currentPath, "~/testdir/../testdir/testdir2")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Test 1 cd path: ", newPath)

}

func TestActionCd2(t *testing.T) {
	const currentPath = "/testdir/testdir2"
	newPath, err := cd(currentPath, "../..")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Test 2 cd path: ", newPath)
}

func TestActionPwd(t *testing.T) {
	const currentPath = "/testdir/testdir2"
	path := pwd(currentPath)
	t.Log(path)
}

func TestActionLs1(t *testing.T) {
	const currentPath = "/"
	output, err := ls(currentPath, []string{"-a"})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(output)
}

func TestActionLs2(t *testing.T) {
	const currentPath = "/"
	output, err := ls(currentPath, []string{})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(output)
}

func TestActionMv1(t *testing.T) {
	const currentPath = "/testdir"
	err := mv(currentPath, "testdir2/", "helloworld/")
	if err != nil {
		t.Fatal(err)
	}
}

func TestActionMv2(t *testing.T) {
	const currentPath = "/testdir"
	err := mv(currentPath, "testfile", "helloworld/")
	if err != nil {
		t.Fatal(err)
	}
}

func TestActionRm1(t *testing.T) {
	const currentPath = "/"
	// Test removing a dir witout flag
	err := rm(currentPath, []string{"testdir"}, []string{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestActionRm2(t *testing.T) {
	const currentPath = "/"
	err := rm(currentPath, []string{"testdir"}, []string{"r", "f"})
	if err != nil {
		t.Fatal(err)
	}
}
