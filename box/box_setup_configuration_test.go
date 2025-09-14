package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCreateConfigurationFolder(t *testing.T) {
	homePath := testDirectory(t)
	defer os.RemoveAll(homePath)

	installedTextBoxWithDummyPackage(t, homePath, "myIdentifier")

	expectedPath := filepath.Join(homePath, "uvbox", "myIdentifier", "configuration")
	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		t.Fatalf("expected folder %s to exist", expectedPath)
	} else if err != nil {
		t.Fatalf("could not check if folder %s exists: %v", expectedPath, err)
	}
}

func TestInstallUvboxConfigurationFile(t *testing.T) {
	homePath := testDirectory(t)
	defer os.RemoveAll(homePath)

	box := installedTextBoxDummy(t, homePath)

	// Check if uvbox configuration file exists
	file, err := box.InstalledUvboxConfigurationFilePath()
	if err != nil {
		t.Fatalf("could not get installed uvbox configuration file path: %v", err)
	}
	if _, err := os.Stat(file); os.IsNotExist(err) {
		t.Fatalf("expected uvbox configuration file %s to exist", file)
	} else if err != nil {
		t.Fatalf("could not check if uvbox configuration file %s exists: %v", file, err)
	}
}

func TestInstalledUvboxConfigurationFileIsUnderConfigurationFolder(t *testing.T) {
	homePath := testDirectory(t)
	defer os.RemoveAll(homePath)

	box := installedTextBoxWithDummyPackage(t, homePath, "myIdentifier")

	file, err := box.InstalledUvboxConfigurationFilePath()
	if err != nil {
		t.Fatalf("could not get installed uvbox configuration file path: %v", err)
	}

	expectedParent := filepath.Join(homePath, "uvbox", "myIdentifier", "configuration")
	_, err = filepath.Rel(expectedParent, file)
	if err != nil {
		t.Fatalf("could not get relative path: %v", err)
	}
}

func TestInstalledUvboxConfigurationFileHasSameContentHasEmbeddedOne(t *testing.T) {
	homePath := testDirectory(t)
	defer os.RemoveAll(homePath)

	box := installedTextBoxDummy(t, homePath)

	installedFile, err := box.InstalledUvboxConfigurationFilePath()
	if err != nil {
		t.Fatalf("could not get installed uvbox configuration file path: %v", err)
	}

	content, err := os.ReadFile(installedFile)
	if err != nil {
		t.Fatalf("could not read file %s: %v", installedFile, err)
	}

	textContent := string(content)
	textEmbeddedContent := string(embeddedConfigurationFile)

	if textContent != textEmbeddedContent {
		t.Fatalf("expected installed content to be equal to embedded content. Got:\n%s\nExpected:\n%s\n", textContent, textEmbeddedContent)
	}
}

func TestInstallCertificatesFile(t *testing.T) {
	homePath := testDirectory(t)
	defer os.RemoveAll(homePath)

	box := installedTextBoxDummy(t, homePath)

	// Check if uvbox configuration file exists
	file, err := box.InstalledCertificatesBundleFilePath()
	if err != nil {
		t.Fatalf("could not get installed certificates bundle file path: %v", err)
	}
	if _, err := os.Stat(file); os.IsNotExist(err) {
		t.Fatalf("expected certificates bundle file %s to exist", file)
	} else if err != nil {
		t.Fatalf("could not check if certificates bundle file %s exists: %v", file, err)
	}
}

func TestInstalledCertificatesBundleFileIsUnderConfigurationFolder(t *testing.T) {
	homePath := testDirectory(t)
	defer os.RemoveAll(homePath)

	box := installedTextBoxWithDummyPackage(t, homePath, "myIdentifier")

	file, err := box.InstalledUvboxConfigurationFilePath()
	if err != nil {
		t.Fatalf("could not get installed certificates bundle file path: %v", err)
	}

	expectedParent := filepath.Join(homePath, "uvbox", "myIdentifier", "configuration")
	_, err = filepath.Rel(expectedParent, file)
	if err != nil {
		t.Fatalf("could not get relative path: %v", err)
	}
}

func TestInstalledCertificatesBundleFileHasSameContentHasEmbeddedOne(t *testing.T) {
	homePath := testDirectory(t)
	defer os.RemoveAll(homePath)

	box := installedTextBoxDummy(t, homePath)

	installedFile, err := box.InstalledCertificatesBundleFilePath()
	if err != nil {
		t.Fatalf("could not get installed uvbox configuration file path: %v", err)
	}

	content, err := os.ReadFile(installedFile)
	if err != nil {
		t.Fatalf("could not read file %s: %v", installedFile, err)
	}

	textContent := string(content)
	textEmbeddedContent := string(embeddedCertificatesBundle)

	if textContent != textEmbeddedContent {
		t.Fatalf("expected installed content to be equal to embedded content. Got:\n%s\nExpected:\n%s\n", textContent, textEmbeddedContent)
	}
}
