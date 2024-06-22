package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
)

func uninstallDevcert() (err error) {
	devcertDir, err := buildDevcertDir()
	if err != nil {
		return
	}

	// Get the devcert executable path
	executable, err := os.Executable()
	if err != nil {
		err = fmt.Errorf("cleaning up CA failed: %w", err)
		return
	}

	// Check for symlinks
	realExecutable, err := filepath.EvalSymlinks(executable)
	if err != nil {
		err = fmt.Errorf("cleaning up CA failed: %w", err)
		return
	}

	color.Set(color.FgWhite, color.Bold)
	fmt.Printf("The uninstall command will:")
	color.Unset()

	fmt.Printf("\n  - Remove the")

	color.Set(color.FgCyan)
	fmt.Printf(" %s ", devcertDir)
	color.Unset()

	fmt.Printf("directory and all the files in it.")

	fmt.Printf("\n  - Remove the")

	color.Set(color.FgCyan)
	fmt.Printf(" %s ", realExecutable)
	color.Unset()

	fmt.Printf("executable.")

	fmt.Printf("\n  - Remove the local devcert Certificate Authority (CA).\n")

	prompt := promptui.Prompt{
		Label:     "Do you want to continue?",
		IsConfirm: true,
	}

	_, err = prompt.Run()

	if err != nil {
		if errors.Is(err, promptui.ErrAbort) {
			err = nil
			return
		}

		err = fmt.Errorf("prompt failed: %w", err)
		return
	}

	// Remove CA from the trust stores
	err = cleanupCA()
	if err != nil {
		err = fmt.Errorf("cleaning up CA failed: %w", err)
		return
	}

	// Remove ~/.devcert
	attemptCleanupDevcertDir()

	// Remove the exeuctable
	os.Remove(realExecutable)

	return
}
