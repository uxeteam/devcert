package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
)

// needsSetup checks if the setup process needs to be executed.
func needsSetup() (need bool, err error) {
	isFolder, err := isDevcertFolder()
	if err != nil {
		return
	}

	if isFolder == false {
		return true, nil
	}

	isCA, err := isValidCA()
	if err != nil {
		return
	}

	if isCA == false {
		return true, nil
	}

	return false, nil
}

// isDevcertFolder checks if the necessary folder exists.
func isDevcertFolder() (is bool, err error) {
	devcertDir, err := buildDevcertDir()
	if err != nil {
		return false, err
	}

	_, err = os.Stat(devcertDir)
	notExist := errors.Is(err, fs.ErrNotExist)

	if err != nil && notExist == false {
		return false, err
	}

	is = notExist == false
	return is, nil
}

// isValidCA check if the certificate authority files are valid.
func isValidCA() (is bool, err error) {
	ca, err := loadCA()
	if err != nil {
		return
	}

	is = ca.Valid

	return
}

// setupPrompt will ask the user to continue or not.
func setupPrompt() (ok bool, err error) {
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

		return
	}

	return true, nil
}

// setup will start the setup process.
func setup() (err error) {
	devcertDir, err := buildDevcertDir()
	if err != nil {
		return
	}

	color.Set(color.FgWhite, color.Bold)
	fmt.Printf("devcert needs to execute the setup process first:")
	color.Unset()

	fmt.Printf("\n  - It will create the")

	color.Set(color.FgCyan)
	fmt.Printf(" %s ", devcertDir)
	color.Unset()

	fmt.Printf("directory.")

	fmt.Printf("\n  - It will create a local certificate authority (CA) to sign future certificates.")
	fmt.Printf("\n  - It will mark the CA as trusted locally.\n")

	promptOK, err := setupPrompt()
	if err != nil {
		return
	}

	if promptOK == false {
		os.Exit(0)
		return
	}

	// Continue the setup process.
	err = createDevcertDir()
	if err != nil {
		err = fmt.Errorf("setup failed: %w", err)
		attemptCleanupDevcertDir()
		return
	}

	err = createCA()
	if err != nil {
		err = fmt.Errorf("setup failed: %w", err)
		attemptCleanupDevcertDir()
		return
	}

	err = trustCA()
	if err != nil {
		err = fmt.Errorf("setup failed: %w", err)
		attemptCleanupCA()
		return
	}

	return
}

// createDevcertDir creates the .devcert directory in the user's home directory.
func createDevcertDir() (err error) {
	fmt.Printf("Creating directory...\n")

	devcertDir, err := buildDevcertDir()
	if err != nil {
		err = fmt.Errorf("creating .devcert directory failed: %w", err)
		return
	}

	isDevcertDir, err := isDevcertFolder()
	if err != nil {
		err = fmt.Errorf("creating .devcert directory failed: %w", err)
		return
	}

	// The directory already exists.
	if isDevcertDir == true {
		fmt.Printf("Directory")

		color.Set(color.FgCyan)
		fmt.Printf(" %s ", devcertDir)
		color.Unset()

		fmt.Printf("is already created.\n")

		return
	}

	// Create the directory
	err = os.MkdirAll(devcertDir, 0755)
	if err != nil {
		err = fmt.Errorf("creating .devcert directory failed: %w", err)
		return
	}

	fmt.Printf("Directory")

	color.Set(color.FgCyan)
	fmt.Printf(" %s ", devcertDir)
	color.Unset()

	fmt.Printf("created.\n")

	return
}
