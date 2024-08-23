package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"syscall"

	"github.com/codeclysm/extract"
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func createTempDir(name string) string {
	var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9 ]+`)
	prefix := nonAlphanumericRegex.ReplaceAllString(name, "_")
	dir, err := os.MkdirTemp("", prefix)
	if err != nil {
		log.Fatal(err)
	}

	return dir
}

func unTar(source string, destination string) error {
	r, err := os.Open(source)
	if err != nil {
		return err
	}
	defer r.Close()

	ctx := context.Background()
	return extract.Archive(ctx, r, destination, nil)
}

func chroot(root string, call string) {
	fmt.Printf("Running %s in %s\n", call, root)
	cmd := exec.Command(call)
	must(syscall.Chroot(root))
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	must(cmd.Run())
}

func pullImage(image string) {
	cmd := exec.Command("./pull", image)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	must(cmd.Run())
}

func main() {
	switch os.Args[1] {
	case "run":
		image := os.Args[2]
		tar := fmt.Sprintf("./assets/%s.tar.gz", image)

		if _, err := os.Stat(tar); errors.Is(err, os.ErrNotExist) {
			panic(err)
		}

		cmd := ""
		if len(os.Args) > 3 {
			cmd = os.Args[3]
		} else {
			buf, err := os.ReadFile(fmt.Sprintf("./assets/%s-cmd", image))
			if err != nil {
				panic(err)
			}
			cmd = string(buf)
		}

		dir := createTempDir(tar)
		defer os.RemoveAll(dir)
		must(unTar(tar, dir))
		chroot(dir, cmd)
	case "pull":
		image := os.Args[2]
		pullImage(image)
	default:
		panic("some error occured")
	}
}
