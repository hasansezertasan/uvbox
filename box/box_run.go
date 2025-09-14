package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func (b *Box) Run(packageName, script string) (int, error) {
	binDir := b.uvToolBinDirPath()
	executablePath := filepath.Join(binDir, script)

	env, err := b.commandsEnvironment()
	if err != nil {
		return 1, fmt.Errorf("could not get uv environment variables: %w", err)
	}

	cmd := exec.Command(executablePath, os.Args[1:]...)
	cmd.Env = env
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err = cmd.Run()
	// Exit with the same exit code as the child process
	if err == nil {
		return 0, nil
	} else if exitError, ok := err.(*exec.ExitError); ok {
		return exitError.ExitCode(), nil
	} else {
		return 1, fmt.Errorf("failed to run command %s: %w", script, err)
	}
}

func (b *Box) RunUv(args []string) (int, error) {
	executable, err := b.InstalledUvExecutablePath()
	if err != nil {
		return 1, fmt.Errorf("could not find uv executable: %w", err)
	}

	env, err := b.commandsEnvironment()
	if err != nil {
		return 1, fmt.Errorf("could not get uv environment variables: %w", err)
	}

	cmd := exec.Command(executable, args...)
	cmd.Env = env
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err = cmd.Run()
	// Exit with the same exit code as the child process
	if err == nil {
		return 0, nil
	} else if exitError, ok := err.(*exec.ExitError); ok {
		return exitError.ExitCode(), nil
	} else {
		return 1, fmt.Errorf("failed to run uv command %s %v: %w", executable, args, err)
	}
}
