package main

import (
	"os"
	"path/filepath"
	"testing"
)

func hideStdoutAndStderr() {
	hideStdout()
	hideStderr()
}

func hideStdout() {
	os.Stdout, _ = os.OpenFile(os.DevNull, os.O_WRONLY, 0)
}

func hideStderr() {
	os.Stderr, _ = os.OpenFile(os.DevNull, os.O_WRONLY, 0)
}

func testDirectory(t *testing.T) string {
	// Create temporary directory
	folder, err := os.MkdirTemp("", "uvbox-test-")
	if err != nil {
		t.Fatalf("could not create temporary directory: %v", err)
	}
	return folder
}

func installedTextBoxWithDummyPackage(t *testing.T, homePath, identifier string) *Box {
	return installedTestBox(t, homePath, identifier, "myPackage")
}

func installedTextBoxDummy(t *testing.T, homePath string) *Box {
	return installedTestBox(t, homePath, "myIdentifier", "myPackage")
}

func installedTestBox(t *testing.T, homePath, identifier string, packageName string) *Box {
	box := GetBox(
		identifier,
		packageName,
		[]string{},
		homePath,
	)

	err := box.Install()
	if err != nil {
		t.Fatal(err)
	}

	return box
}

func TestCreateUvboxFolder(t *testing.T) {
	homePath := testDirectory(t)
	defer os.RemoveAll(homePath)

	installedTextBoxDummy(t, homePath)

	expectedPath := filepath.Join(homePath, "uvbox")
	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		t.Fatalf("expected folder %s to exist", expectedPath)
	} else if err != nil {
		t.Fatalf("could not check if folder %s exists: %v", expectedPath, err)
	}
}

func TestCanTellIfNeedsInstall(t *testing.T) {
	homePath := testDirectory(t)
	defer os.RemoveAll(homePath)

	box := GetBox("myIdentifier", "myPackage", []string{}, homePath)

	if yes, err := box.NeedsInstall(); err != nil {
		t.Fatalf("could not check if needs install: %v", err)
	} else if !yes {
		t.Fatalf("expected to need install")
	}
}

func TestCanTellIfDoesNotNeedInstall(t *testing.T) {
	homePath := testDirectory(t)
	defer os.RemoveAll(homePath)

	box := installedTextBoxDummy(t, homePath)

	if yes, err := box.NeedsInstall(); err != nil {
		t.Fatalf("could not check if needs install: %v", err)
	} else if yes {
		t.Fatalf("expected to not need install")
	}
}

func TestCanInstallWithSpinner(t *testing.T) {
	hideStdoutAndStderr()

	homePath := testDirectory(t)
	defer os.RemoveAll(homePath)

	box := GetBox("myIdentifier", "myPackage", []string{}, homePath)

	err := box.InstallWithSpinner()
	if err != nil {
		t.Fatalf("could not install: %v", err)
	}

	if yes, err := box.NeedsInstall(); err != nil {
		t.Fatalf("could not check if needs install: %v", err)
	} else if yes {
		t.Fatalf("expected to not need install")
	}
}

func TestCanUninstall(t *testing.T) {
	hideStdoutAndStderr()

	homePath := testDirectory(t)
	defer os.RemoveAll(homePath)

	box := installedTextBoxWithPackage(t, homePath, "cowsay", "6.1", "")

	err := box.Uninstall()
	if err != nil {
		t.Fatalf("could not uninstall box: %v", err)
	}

	if yes, err := box.uvboxConfigurationInstalled(); err != nil {
		t.Fatalf("could not check if needs configuration needs installation: %v", err)
	} else if yes {
		t.Fatalf("expected configuration file to need installation after uninstallation")
	}

	if yes, err := box.certificatesBundleInstalled(); err != nil {
		t.Fatalf("could not check if needs certificates bundle file needs installation: %v", err)
	} else if yes {
		t.Fatalf("expected certificates bundle file to need installation after uninstallation")
	}

	if yes, err := box.IsPackageInstalled(); err != nil {
		t.Fatalf("could not check if package needs installation: %v", err)
	} else if yes {
		t.Fatalf("expected package to need installation after uninstallation")
	}
}

func TestCanUninstallWithoutPriorInstallWithoutCrash(t *testing.T) {
	hideStdoutAndStderr()

	homePath := testDirectory(t)
	defer os.RemoveAll(homePath)

	box := GetBox("myIdentifier", "myPackage", []string{}, homePath)

	err := box.Uninstall()
	if err != nil {
		t.Fatalf("could not uninstall box: %v", err)
	}
}
