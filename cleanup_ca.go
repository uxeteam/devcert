package main

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

// cleanupCA removes the trusted CA certificate from the OS key store.
func cleanupCA() (err error) {
	fmt.Printf("Removing certificate authority configuration...\n")

	ca, err := loadCA()
	if err != nil {
		err = fmt.Errorf("removing CA configuration failed: %w", err)
		return
	}

	if ca.Valid == false {
		err = fmt.Errorf("removing CA configuration failed: certificate is invalid")
		return
	}

	crtPath, _, err := buildCAPaths()
	if err != nil {
		err = fmt.Errorf("removing CA configuration failed: %w", err)
		return
	}

	switch runtime.GOOS {
	case "darwin":
		err = deleteCADarwin()
	case "windows":
		err = deleteCAWindows()
	case "linux":
		err = deleteCALinux(crtPath)
	}

	if err == nil {
		fmt.Printf("Certificate authority (CA) is uninstalled.\n")
	}

	return
}

func deleteCADarwin() (err error) {
	// Remove trust and delete the certificate from Keychain
	std, err := sudoify("security", "delete-certificate", "-c", commonNameCA, "-t").CombinedOutput()
	if err != nil {
		err = fmt.Errorf("deleting CA failed: %s, %w", std, err)
		return
	}

	return
}

func deleteCAWindows() (err error) {
	var argList strings.Builder
	argList.WriteString("-ArgumentList '-delstore -f ROOT ''")
	argList.WriteString(commonNameCA)
	argList.WriteString("'''")

	stdOutStdError, err := exec.Command("powershell", "Start-Process -FilePath certutil -Verb RunAs -Wait -PassThru", argList.String()).CombinedOutput()
	if err != nil {
		err = fmt.Errorf("Trusting CA failed: %s, %w", stdOutStdError, err)
		return
	}

	return
}

func deleteCALinux(crtPath string) (err error) {
	base, err := detectLinux()
	if err != nil {
		err = fmt.Errorf("deleting CA failed: %w", err)
		return
	}

	switch base {
	case "debian":
		stdOutStdError, err := exec.Command("rm", "-rf", "/usr/local/share/ca-certificates/devcert_ca.crt").CombinedOutput()
		if err != nil {
			err = fmt.Errorf("deleting CA failed: %s, %w", stdOutStdError, err)
			return err
		}

		stdOutStdError, err = sudoify("update-ca-certificates", "--fresh").CombinedOutput()
		if err != nil {
			err = fmt.Errorf("deleting CA failed: %s, %w", stdOutStdError, err)
			return err
		}

	case "rhel":
		stdOutStdError, err := exec.Command("rm", "-rf", "/etc/pki/ca-trust/source/anchors/devcert_ca.crt").CombinedOutput()
		if err != nil {
			err = fmt.Errorf("deleting CA failed: %s, %w", stdOutStdError, err)
			return err
		}

		stdOutStdError, err = sudoify("update-ca-trust", "extract").CombinedOutput()
		if err != nil {
			err = fmt.Errorf("deleting CA failed: %s, %w", stdOutStdError, err)
			return err
		}

	case "arch":
		stdOutStdError, err := exec.Command("trust", "anchor", "--remove", crtPath).CombinedOutput()
		if err != nil {
			err = fmt.Errorf("deleting CA failed: %s, %w", stdOutStdError, err)
			return err
		}
	}

	return
}
