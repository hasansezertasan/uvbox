package main

import (
	"os"
	"testing"
)

func TestForwardArgumentsAndReturnCode(t *testing.T) {
	hideStdout()

	homePath := testDirectory(t)
	defer os.RemoveAll(homePath)

	packageToInstall := "cowsay"
	packageVersion := "6.1"
	box := installedTextBoxWithPackage(t, homePath, packageToInstall, packageVersion, "")

	// cowsay --help returns code 0
	os.Args = []string{"cowsay", "--help"}

	expectedCode := 0
	returnCode, err := box.Run("cowsay", "cowsay")
	if err != nil {
		t.Fatalf("could not run package %s: %v", packageToInstall, err)
	}

	if returnCode != expectedCode {
		t.Fatalf("expected exit code %d, got %d", expectedCode, returnCode)
	}
}

func TestCanRunInstalledPackageAndForwardReturnFailureCode(t *testing.T) {
	hideStdoutAndStderr()

	homePath := testDirectory(t)
	defer os.RemoveAll(homePath)

	packageToInstall := "cowsay"
	packageVersion := "6.1"
	box := installedTextBoxWithPackage(t, homePath, packageToInstall, packageVersion, "")

	os.Args = []string{"cowsay", "--bad-argument"}

	expectedCode := 2
	returnCode, err := box.Run("cowsay", "cowsay")
	if err != nil {
		t.Fatalf("could not run package %s: %v", packageToInstall, err)
	}

	if returnCode != expectedCode {
		t.Fatalf("expected exit code %d, got %d", expectedCode, returnCode)
	}
}
