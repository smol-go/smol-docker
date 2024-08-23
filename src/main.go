package main

import (
	"os"
	"os/exec"
	"syscall"
)

func main() {
	cmd := exec.Command("/hello")
	syscall.Chroot("/testroot")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}
