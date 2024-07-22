package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	runCmd := flag.NewFlagSet("run", flag.ExitOnError)
	buildCmd := flag.NewFlagSet("build", flag.ExitOnError)

	if len(os.Args) <= 1 {
		fmt.Println("please enter a valid subcommand.")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "run":
		runCmd.Parse(os.Args[2:])
		arg := runCmd.Args()
		run(arg)
	case "child":
		runCmd.Parse(os.Args[2:])
		arg := runCmd.Args()
		child(arg)
	case "build":
		tag := buildCmd.String("tag", "", "Name of container image")
		path := buildCmd.String("path", "", "Path to ContainerFile")
		buildCmd.Parse(os.Args[2:])
		build(*tag, *path)
	default:
		fmt.Printf("invalid subcommand %s", os.Args[1])
		os.Exit(1)
	}
}

func run(args []string) {
	fmt.Printf("Running %v \n", args)
}

func child(args []string) {
	fmt.Printf("Running from proc in namespace %v \n", args)
}

func build(tag, path string) {
	fmt.Printf("running build with tag: %s and %s", tag, path)
}
