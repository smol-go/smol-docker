package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"time"

	"github.com/codeclysm/extract"
)

const (
	dumpsDir       = "./dumps"
	defaultTimeout = 5 * time.Minute
	pullScript     = "./pull.sh"
)

type Container struct {
	Image   string
	Command string
	RootDir string
	TempDir string
}

func NewContainer(image, command string) *Container {
	return &Container{
		Image:   image,
		Command: command,
		RootDir: filepath.Join(dumpsDir, image),
	}
}

func (c *Container) Setup() error {
	tempDir, err := c.createTempDir()
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	c.TempDir = tempDir

	tarPath := filepath.Join(c.RootDir, fmt.Sprintf("%s.tar.gz", c.Image))
	if err := c.unpackImage(tarPath); err != nil {
		return fmt.Errorf("failed to unpack image: %w", err)
	}

	return nil
}

func (c *Container) Run() error {
	defer c.Cleanup()

	if runtime.GOOS == "linux" {
		if err := c.setupMounts(); err != nil {
			return fmt.Errorf("failed to setup mounts: %w", err)
		}

		if err := c.changeRoot(); err != nil {
			return fmt.Errorf("failed to chroot: %w", err)
		}
	} else {
		// macOS implementation
		if err := c.runMacOS(); err != nil {
			// If the error already contains a formatted message about Linux binaries,
			// return it directly without wrapping
			if err.Error() == "cannot run Linux binaries on macOS. Please use Docker Desktop or a Linux VM to run containers" {
				return err
			}
			return fmt.Errorf("failed to run container on macOS: %w", err)
		}
	}

	return nil
}

func (c *Container) Cleanup() {
	if c.TempDir != "" {
		if err := os.RemoveAll(c.TempDir); err != nil {
			log.Printf("Warning: failed to cleanup temp directory: %v", err)
		}
	}
}

func (c *Container) createTempDir() (string, error) {
	nonAlphanumericRegex := regexp.MustCompile(`[^a-zA-Z0-9 ]+`)
	prefix := nonAlphanumericRegex.ReplaceAllString(c.Image, "_")

	dir, err := os.MkdirTemp("", prefix)
	if err != nil {
		return "", fmt.Errorf("failed to create temp directory: %w", err)
	}

	return dir, nil
}

func (c *Container) unpackImage(tarPath string) error {
	if _, err := os.Stat(tarPath); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("tar file not found: %s", tarPath)
	}

	r, err := os.Open(tarPath)
	if err != nil {
		return fmt.Errorf("failed to open tar: %w", err)
	}
	defer r.Close()

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	if err := extract.Archive(ctx, r, c.TempDir, nil); err != nil {
		return fmt.Errorf("failed to extract archive: %w", err)
	}

	return nil
}

func (c *Container) setupMounts() error {
	// This is a stub that will be overridden by the platform-specific implementation
	return nil
}

func (c *Container) changeRoot() error {
	// This is a stub that will be overridden by the platform-specific implementation
	return nil
}

func (c *Container) runMacOS() error {
	// For macOS, we'll use a simpler approach that doesn't require root privileges
	// We'll just run the command in the extracted directory with modified environment
	
	// Make the command path relative to the extracted directory
	cmdPath := c.Command
	if filepath.IsAbs(cmdPath) {
		cmdPath = filepath.Join(c.TempDir, cmdPath[1:])
	} else {
		cmdPath = filepath.Join(c.TempDir, cmdPath)
	}
	
	// Check if the file exists and is executable
	if _, err := os.Stat(cmdPath); os.IsNotExist(err) {
		return fmt.Errorf("executable not found in container: %s", cmdPath)
	}
	
	// Try to execute the command
	cmd := exec.Command(cmdPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	// Set the working directory to the extracted image
	cmd.Dir = c.TempDir
	
	// Set environment variables to simulate container environment
	cmd.Env = append(os.Environ(),
		"CONTAINER=true",
		"HOME="+filepath.Join(c.TempDir, "home"),
		"PATH="+filepath.Join(c.TempDir, "bin")+":"+os.Getenv("PATH"),
	)

	if err := cmd.Run(); err != nil {
		if execErr, ok := err.(*exec.Error); ok && execErr.Err == exec.ErrNotFound {
			return fmt.Errorf("executable not found: %s", cmdPath)
		}
		if execErr, ok := err.(*exec.ExitError); ok {
			return fmt.Errorf("command failed with exit code %d: %s", execErr.ExitCode(), execErr.String())
		}
		if err.Error() == "exec format error" {
			return fmt.Errorf("cannot run Linux binaries on macOS. Please use Docker Desktop or a Linux VM to run containers")
		}
		return fmt.Errorf("command failed: %w", err)
	}

	return nil
}

func pullImage(image string) error {
	cmd := exec.Command(pullScript, image)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func parseArgs() (string, string, string, error) {
	if len(os.Args) < 2 {
		return "", "", "", fmt.Errorf("insufficient arguments")
	}

	command := os.Args[1]

	if command != "run" && command != "pull" {
		return "", "", "", fmt.Errorf("invalid command: %s", command)
	}

	if len(os.Args) < 3 {
		return "", "", "", fmt.Errorf("image name required")
	}

	image := os.Args[2]
	var cmd string

	if command == "run" {
		if len(os.Args) > 3 {
			cmd = os.Args[3]
		} else {
			cmdFile := filepath.Join(dumpsDir, image, fmt.Sprintf("%s-cmd", image))
			buf, err := os.ReadFile(cmdFile)
			if err != nil {
				return "", "", "", fmt.Errorf("failed to read command file: %w", err)
			}
			cmd = string(buf)
		}
	}

	return command, image, cmd, nil
}

func main() {
	command, image, cmd, err := parseArgs()
	if err != nil {
		log.Fatalf("Failed to parse arguments: %v", err)
	}

	switch command {
	case "run":
		if err := pullImage(image); err != nil {
			log.Fatalf("Failed to pull image: %v", err)
		}

		container := NewContainer(image, cmd)

		if err := container.Setup(); err != nil {
			log.Fatalf("Failed to setup container: %v", err)
		}

		if err := container.Run(); err != nil {
			log.Fatalf("Failed to run container: %v", err)
		}

	case "pull":
		if err := pullImage(image); err != nil {
			log.Fatalf("Failed to pull image: %v", err)
		}
	}
}
