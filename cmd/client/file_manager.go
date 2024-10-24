package main

import (
	"fmt"
)

// FileManger gets the current directory from the server
type FileManger struct {
	CurrentDir string
	Flags      []string
}

// Commands send the tasks to the server
func (fm *FileManger) Pwd() {
	fmt.Printf("%v", fm.CurrentDir)
}

func (fm *FileManger) Mkdir(paths ...string) {
	fmt.Println("mkdir")
}

func (fm *FileManger) Ls(paths ...string) {
	fmt.Println("ls")
}

func (fm *FileManger) Cd(paths ...string) {
	fmt.Println("Cd")
}

func (fm *FileManger) Rm(paths ...string) {
	fmt.Println("Rm")
}

// Upload send all the files to the server
func (fm *FileManger) Upload(paths ...string) {
	fmt.Println("Upload")
}

// Download works for every selected file
func (fm *FileManger) Download(paths ...string) {
	fmt.Println("Download")
}
