//go:build linux
// +build linux

package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func (c *Container) setupMounts() error {
	procPath := filepath.Join(c.TempDir, "proc")
	if err := os.MkdirAll(procPath, 0755); err != nil {
		return fmt.Errorf("failed to create proc directory: %w", err)
	}

	if err := syscall.Mount("proc", procPath, "proc", 0, ""); err != nil {
		return fmt.Errorf("failed to mount proc: %w", err)
	}

	return nil
}

func (c *Container) changeRoot() error {
	oldRoot, err := os.Open("/")
	if err != nil {
		return fmt.Errorf("failed to open root: %w", err)
	}
	defer oldRoot.Close()

	cmd := exec.Command(c.Command)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := syscall.Chdir(c.TempDir); err != nil {
		return fmt.Errorf("failed to change directory: %w", err)
	}

	if err := syscall.Chroot(c.TempDir); err != nil {
		return fmt.Errorf("failed to chroot: %w", err)
	}

	if err := cmd.Run(); err != nil {
		log.Printf("Command failed: %v", err)
	}

	if err := syscall.Fchdir(int(oldRoot.Fd())); err != nil {
		return fmt.Errorf("failed to restore old root directory: %w", err)
	}

	if err := syscall.Chroot("."); err != nil {
		return fmt.Errorf("failed to restore old root: %w", err)
	}

	return nil
} 