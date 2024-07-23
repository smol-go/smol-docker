package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
)

type limit struct {
	name  string
	path  string
	param string
	value []byte
}

type limits struct {
	Limits []limit
}

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

	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID,
	}

	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}

func child(args []string) {
	fmt.Printf("Running from proc in namespace %v \n", args)

	err := cgroup()
	if err != nil {
		panic(err)
	}

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = syscall.Sethostname([]byte("container"))
	if err != nil {
		panic(err)
	}

	if err = syscall.Chroot("ubuntu-rootfs/"); err != nil {
		panic(err)
	}

	if err = syscall.Chdir("/"); err != nil {
		panic(err)
	}

	if err = syscall.Mount("proc", "proc", "proc", 0, ""); err != nil {
		panic(err)
	}

	err = cmd.Run()
	if err != nil {
		panic(err)
	}

	if err = syscall.Unmount("proc", 0); err != nil {
		panic(err)
	}
}

func build(tag, path string) {
	fmt.Printf("running build with tag: %s and %s", tag, path)
}

func cgroup() error {
	cgrouplimits := limits{
		[]limit{
			{
				"pids",
				"/sys/fs/cgroup/pids/container",
				"pids.max",
				[]byte("20"),
			},
			{
				"memory",
				"/sys/fs/cgroup/memory/container",
				"memory.limit_in_bytes",
				[]byte("1000000"),
			},
			{
				"cpu",
				"/sys/fs/cgroup/cpu/container",
				"cpu.shares",
				[]byte("512"),
			},
		},
	}

	for _, l := range cgrouplimits.Limits {
		fmt.Println(filepath.Join(l.path, l.param))
		//os.Mkdir
		os.Mkdir(l.path, 0755)
		//Create cgroup limit
		err := os.WriteFile(filepath.Join(l.path, l.param), l.value, 0700)
		if err != nil {
			return err
		}
		//os.WriteFile (add proc to cgroup)
		err = os.WriteFile(filepath.Join(l.path, "cgroup.procs"), []byte(strconv.Itoa(os.Getpid())), 0700)
		if err != nil {
			return err
		}
		//os.WriteFile notify_on-release
		err = os.WriteFile(filepath.Join(l.path, "notify_on_release"), []byte("1"), 0700)
		if err != nil {
			return err
		}
	}

	return nil
}
