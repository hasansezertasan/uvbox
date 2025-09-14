package main

import (
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func testDirectory(t *testing.T) string {
	folder, err := os.MkdirTemp("", "uvbox-boxer-test")
	if err != nil {
		t.Fatalf("could not create temporary directory: %v", err)
	}
	return folder
}

func testDirectoryWithFile(t *testing.T, filename string) (string, string) {
	workingDirectory := testDirectory(t)

	// Create file
	filePath := filepath.Join(workingDirectory, filename)
	if _, err := os.Create(filePath); err != nil {
		os.RemoveAll(workingDirectory)
		t.Fatalf("could not write pyproject.toml file under test temporary directory at: %s", filePath)
	}

	return workingDirectory, filePath
}

func testDirectoryWithFileContent(t *testing.T, filename, content string) (string, string) {
	workingDirectory, filePath := testDirectoryWithFile(t, filename)

	// Write file content
	if err := os.WriteFile(filePath, []byte(content), 0700); err != nil {
		t.Errorf("error while writing pyproject.toml file content at %s: %s", filePath, err)
	}

	return workingDirectory, filePath
}

func testChildDirectoryWithFileContentInParent(t *testing.T, filename, content string) (string, string) {
	// Parent
	workingDirectory, filePath := testDirectoryWithFileContent(t, filename, content)

	// Create child directory
	child := filepath.Join(workingDirectory, "child")
	if err := os.Mkdir(child, 0700); err != nil {
		t.Errorf("failed to create child test directoryat %s: %s", child, err)
	}

	return child, filePath
}

func testCanDetectFileInDirectory(t *testing.T, filename string) {
	workingDirectory, _ := testDirectoryWithFile(t, filename)
	defer os.RemoveAll(workingDirectory)

	if foundPath, err := directoryOrParentsContainsFile(workingDirectory, filename); err != nil {
		t.Errorf("error while checking if test working directory contains %s file: %s", filename, err)
	} else if foundPath == "" {
		t.Fatalf("directory is supposed to contain a %s file: %s", filename, workingDirectory)
	}
}

func testCanDetectFileInParentDirectory(t *testing.T, fileName string) {
	workingDirectory := testDirectory(t)
	defer os.RemoveAll(workingDirectory)

	childDirectory := filepath.Join(workingDirectory, "child")

	// Create a child directory
	if err := os.MkdirAll(childDirectory, 0700); err != nil {
		t.Errorf("could not create child directory under test directory at %s: %s", childDirectory, err)
	}

	// Create file at parent directory
	filePath := filepath.Join(workingDirectory, fileName)
	if _, err := os.Create(filePath); err != nil {
		t.Errorf("Error while creating %s under child directory at %s: %s", fileName, childDirectory, err)
	}

	// Test file is containing with a call with the child directory
	if foundPath, err := directoryOrParentsContainsFile(childDirectory, fileName); err != nil {
		t.Errorf("error while checking if test working directory contains %s file: %s", fileName, err)
	} else if foundPath == "" {
		t.Fatalf("directory, or its parents, are supposed to contain a %s file: %s", fileName, workingDirectory)
	}
}

func TestCanDetectPyProjectTomlInParentDirectories(t *testing.T) {
	testCanDetectFileInParentDirectory(t, "pyproject.toml")
}

func TestCanDetectUvboxTomlInParentDirectories(t *testing.T) {
	testCanDetectFileInParentDirectory(t, "uvbox.toml")
}

func testCanDetectFileMissingInDirectory(t *testing.T, filename string) {
	workingDirectory := testDirectory(t)
	defer os.RemoveAll(workingDirectory)

	if foundPath, err := directoryOrParentsContainsFile(workingDirectory, filename); err != nil {
		t.Errorf("error while checking if test working directory contains %s file: %s", filename, err)
	} else if foundPath != "" {
		t.Fatalf("directory is not supposed to contain a %s file: %s", filename, workingDirectory)
	}
}

func TestCanDetectIfUvboxTomlInDirectory(t *testing.T) {
	testCanDetectFileInDirectory(t, "uvbox.toml")
}

func TestCanDetectIfPyprojectTomlInDirectory(t *testing.T) {
	testCanDetectFileInDirectory(t, "pyproject.toml")
}

func TestCanDetectUvboxTomlMissingInDirectory(t *testing.T) {
	testCanDetectFileMissingInDirectory(t, "uvbox.toml")
}

func TestCanDetectPyprojectTomlMissingInDirectory(t *testing.T) {
	testCanDetectFileMissingInDirectory(t, "pyproject.toml")
}

func TestCanDetectIfPyprojectTomlMissesToolUvbox(t *testing.T) {
	content := `
		[project]
		name = "your_project_name"
		version = "0.1.0"
		description = "A short description of your project."
		authors = [
		    {name = "Your Name", email = "your.email@example.com"}
		]
		dependencies = [
		    "package_name>=1.0",
		]

		[build-system]
		requires = ["setuptools>=61.0"]
		build-backend = "setuptools.build_meta"
	`
	workingDirectory, pyprojectFile := testDirectoryWithFileContent(t, "pyproject.toml", content)
	defer os.RemoveAll(workingDirectory)

	if ok, err := tomlContainsToolUvbox(pyprojectFile); err != nil {
		t.Errorf("could not check if test toml file contains tool.uvbox: %s", err)
	} else if ok {
		t.Fatalf("tool.uvbox should not have been detected inside test toml file: %s", pyprojectFile)
	}
}

func TestCanDetectIfPyprojectTomlContainsToolUvbox(t *testing.T) {
	content := `
		[project]
		name = "your_project_name"
		version = "0.1.0"
		description = "A short description of your project."
		authors = [
		    {name = "Your Name", email = "your.email@example.com"}
		]
		dependencies = [
		    "package_name>=1.0",
		]

		[build-system]
		requires = ["setuptools>=61.0"]
		build-backend = "setuptools.build_meta"

		[tool.uvbox]
		something = "ok"
	`
	workingDirectory, pyprojectFile := testDirectoryWithFileContent(t, "pyproject.toml", content)
	defer os.RemoveAll(workingDirectory)

	if ok, err := tomlContainsToolUvbox(pyprojectFile); err != nil {
		t.Errorf("could not check if test toml file contains tool.uvbox: %s", err)
	} else if !ok {
		t.Fatalf("tool.uvbox should be detected inside test toml file: %s", pyprojectFile)
	}
}

func TestCanDetectIfPyprojectTomlContainsToolUvboxWhenOnlyUvboxPackageName(t *testing.T) {
	content := `
		[project]
		name = "your_project_name"
		version = "0.1.0"
		description = "A short description of your project."
		authors = [
		    {name = "Your Name", email = "your.email@example.com"}
		]
		dependencies = [
		    "package_name>=1.0",
		]

		[build-system]
		requires = ["setuptools>=61.0"]
		build-backend = "setuptools.build_meta"

		[tool.uvbox.package]
		name = "my-package"
	`
	workingDirectory, pyprojectFile := testDirectoryWithFileContent(t, "pyproject.toml", content)
	defer os.RemoveAll(workingDirectory)

	if ok, err := tomlContainsToolUvbox(pyprojectFile); err != nil {
		t.Errorf("could not check if test toml file contains tool.uvbox: %s", err)
	} else if !ok {
		t.Fatalf("tool.uvbox should be detected inside test toml file: %s", pyprojectFile)
	}
}

func TestCanReadConfigPackageFromPyprojectToml(t *testing.T) {
	content := `
		[project]
		name = "your_project_name"
		version = "0.1.0"
		description = "A short description of your project."
		authors = [
		    {name = "Your Name", email = "your.email@example.com"}
		]
		dependencies = [
		    "package_name>=1.0",
		]

		[build-system]
		requires = ["setuptools>=61.0"]
		build-backend = "setuptools.build_meta"

		[tool.uvbox.package]
		name = "my-package"
		script = "my.script:run"
	`
	workingDirectory, pyprojectFile := testDirectoryWithFileContent(t, "pyproject.toml", content)
	defer os.RemoveAll(workingDirectory)

	// Parse tool.uvbox.package
	configuration, err := loadConfigurationFromPyprojectToml(pyprojectFile)
	if err != nil {
		t.Errorf("error while parsing uvbox package configuration from pyproject.toml: %s", err)
	}

	packageConfiguration, err := configuration.ParsedPackageConfiguration()
	if err != nil {
		t.Errorf("error while parsing package configuration: %s", err)
	}

	// Check values
	packageNameValue := packageConfiguration.Name
	packageScriptValue := packageConfiguration.Script
	if packageNameValue != "my-package" {
		t.Errorf("decoded package name value should have been 'my-package', current is: %s", packageNameValue)
	}
	if packageScriptValue != "my.script:run" {
		t.Errorf("decoded package script value should bave been 'my.script:run', current is: %s", packageScriptValue)
	}
}

func assertPackageNameAndScriptFromAutomaticallyFoundConfiguration(t *testing.T, workingDirectory, expectedPackageName, expectedPackageScriptName string) {
	configuration, err := loadConfigurationFromEitherCLIArgumentOrDetectedFile(workingDirectory)
	if err != nil {
		t.Errorf("could not load test configuration: %s", err)
	}
	assertPackageNameAndScriptFromConfiguration(t, configuration, expectedPackageName, expectedPackageScriptName)
}

func assertPackageNameAndScriptFromConfiguration(t *testing.T, configuration Configuration, expectedPackageName, expectedPackageScriptName string) {
	// Parse tool.uvbox.package
	packageConfiguration, err := configuration.ParsedPackageConfiguration()
	if err != nil {
		t.Errorf("error while parsing package configuration: %s", err)
	}

	// Check values
	packageNameValue := packageConfiguration.Name
	packageScriptValue := packageConfiguration.Script
	if packageNameValue != expectedPackageName {
		t.Errorf("decoded package name value should have been '%s', current is: %s", expectedPackageName, packageNameValue)
	}
	if packageScriptValue != "my.script:run" {
		t.Errorf("decoded package script value should bave been '%s', current is: %s", expectedPackageScriptName, packageScriptValue)
	}
}

func TestCanReadConfigPackageFromUvboxToml(t *testing.T) {
	content := `
		[package]
		name = "my-package"
		script = "my.script:run"
	`

	// Directory with test configuration file
	workingDirectory, uvboxFile := testDirectoryWithFileContent(t, "uvbox.toml", content)
	defer os.RemoveAll(workingDirectory)

	// Parse uvbox.toml
	configuration, err := loadConfigurationFromUvboxToml(uvboxFile)
	if err != nil {
		t.Errorf("error while loading test pyproject.toml located at %s: %s", uvboxFile, err)
	}

	assertPackageNameAndScriptFromConfiguration(t, configuration, "my-package", "my.script:run")
}

func TestCanLoadUvboxTomlFromCLIArgumentValue(t *testing.T) {
	content := `
		[package]
		name = "my-package"
		script = "my.script:run"
	`

	// Directory with test configuration file
	workingDirectory, uvboxFile := testDirectoryWithFileContent(t, "uvbox.toml", content)
	defer os.RemoveAll(workingDirectory)

	// Simulate CLI configuration file argument
	Config = uvboxFile

	// Confirm parsed configuration is valid
	assertPackageNameAndScriptFromAutomaticallyFoundConfiguration(t, workingDirectory, "my-package", "my.script:run")
}

func TestCanLoadUvboxTomlFromWorkingDirectory(t *testing.T) {
	content := `
		[package]
		name = "my-package"
		script = "my.script:run"
	`

	// Directory with test configuration file
	workingDirectory, _ := testDirectoryWithFileContent(t, "uvbox.toml", content)
	defer os.RemoveAll(workingDirectory)

	// Disable CLI configuration file argument
	Config = ""

	// Confirm parsed configuration is valid
	assertPackageNameAndScriptFromAutomaticallyFoundConfiguration(t, workingDirectory, "my-package", "my.script:run")
}

func TestCanLoadUvboxTomlFromWorkingDirectoryParent(t *testing.T) {
	content := `
		[package]
		name = "my-package"
		script = "my.script:run"
	`

	// Directory with test configuration file
	workingDirectory, _ := testChildDirectoryWithFileContentInParent(t, "uvbox.toml", content)
	defer os.RemoveAll(workingDirectory)

	// Disable CLI configuration file argument
	Config = ""

	// Confirm parsed configuration is valid
	assertPackageNameAndScriptFromAutomaticallyFoundConfiguration(t, workingDirectory, "my-package", "my.script:run")
}

func TestCanLoadPyprojectTomlFromCLIArgumentValue(t *testing.T) {
	content := `
		[project]
		name = "your_project_name"
		version = "0.1.0"
		description = "A short description of your project."
		authors = [
		    {name = "Your Name", email = "your.email@example.com"}
		]
		dependencies = [
		    "package_name>=1.0",
		]

		[build-system]
		requires = ["setuptools>=61.0"]
		build-backend = "setuptools.build_meta"

		[tool.uvbox.package]
		name = "my-package-from-pyproject"
		script = "my.script:run"
	`

	// Directory with test configuration file
	workingDirectory, pyprojectFile := testDirectoryWithFileContent(t, "pyproject.toml", content)
	defer os.RemoveAll(workingDirectory)

	// Simulate CLI configuration file argument
	Config = pyprojectFile

	// Confirm parsed configuration is valid
	assertPackageNameAndScriptFromAutomaticallyFoundConfiguration(t, workingDirectory, "my-package-from-pyproject", "my.script:run")
}

func TestCanLoadPyprojectTomlFromWorkingDirectory(t *testing.T) {
	content := `
		[project]
		name = "your_project_name"
		version = "0.1.0"
		description = "A short description of your project."
		authors = [
		    {name = "Your Name", email = "your.email@example.com"}
		]
		dependencies = [
		    "package_name>=1.0",
		]

		[build-system]
		requires = ["setuptools>=61.0"]
		build-backend = "setuptools.build_meta"

		[tool.uvbox.package]
		name = "my-package-from-pyproject"
		script = "my.script:run"
	`

	// Directory with test configuration file
	workingDirectory, _ := testDirectoryWithFileContent(t, "pyproject.toml", content)
	defer os.RemoveAll(workingDirectory)

	// Disable CLI configuration file argument
	Config = ""

	// Confirm parsed configuration is valid
	assertPackageNameAndScriptFromAutomaticallyFoundConfiguration(t, workingDirectory, "my-package-from-pyproject", "my.script:run")
}

func TestCanLoadPyprojectTomlFromWorkingDirectoryParent(t *testing.T) {
	content := `
		[project]
		name = "your_project_name"
		version = "0.1.0"
		description = "A short description of your project."
		authors = [
		    {name = "Your Name", email = "your.email@example.com"}
		]
		dependencies = [
		    "package_name>=1.0",
		]

		[build-system]
		requires = ["setuptools>=61.0"]
		build-backend = "setuptools.build_meta"

		[tool.uvbox.package]
		name = "my-package-from-pyproject"
		script = "my.script:run"
	`

	// Directory with test configuration file
	workingDirectory, _ := testChildDirectoryWithFileContentInParent(t, "pyproject.toml", content)
	defer os.RemoveAll(workingDirectory)

	// Disable CLI configuration file argument
	Config = ""

	// Confirm parsed configuration is valid
	assertPackageNameAndScriptFromAutomaticallyFoundConfiguration(t, workingDirectory, "my-package-from-pyproject", "my.script:run")
}

func TestSaveConfigurationToDirectory(t *testing.T) {
	content := `
		[package]
		name = "my-package"
		script = "my.script:run"
	`

	// Directory with test configuration file
	workingDirectory, filePath := testDirectoryWithFileContent(t, "test.toml", content)
	defer os.RemoveAll(workingDirectory)

	// Disable CLI configuration file argument
	Config = filePath

	// Save configuration
	configuration, err := loadConfigurationFromEitherCLIArgumentOrDetectedFile(workingDirectory)
	if err != nil {
		t.Errorf("could not load test configuration: %s", err)
	}
	saveConfigurationToDirectory(configuration, workingDirectory, "copy.toml")

	// Load it back
	Config = filepath.Join(workingDirectory, "copy.toml")
	assertPackageNameAndScriptFromAutomaticallyFoundConfiguration(t, workingDirectory, "my-package", "my.script:run")
}

func TestAppendUvPythonInstallMirror(t *testing.T) {
	content := `
		[project]
		name = "your_project_name"
		version = "0.1.0"
		description = "A short description of your project."
		authors = [
		    {name = "Your Name", email = "your.email@example.com"}
		]
		dependencies = [
		    "package_name>=1.0",
		]

		[build-system]
		requires = ["setuptools>=61.0"]
		build-backend = "setuptools.build_meta"

		[tool.uv]
		python-install-mirror = "https://mirror.com"

		[tool.uvbox.package]
		name = "my-package-from-pyproject"
		script = "my.script:run"
	`

	// Directory with test configuration file
	workingDirectory, _ := testDirectoryWithFileContent(t, "pyproject.toml", content)
	defer os.RemoveAll(workingDirectory)

	// Disable CLI configuration file argument
	Config = ""

	assertUvEnvironmentVariableInAutomaticallyFoundConfiguration(t, workingDirectory, "UV_PYTHON_INSTALL_MIRROR=https://mirror.com")
}

func TestAppendUvDefaultIndexUnnamed(t *testing.T) {
	content := `
		[project]
		name = "your_project_name"
		version = "0.1.0"
		description = "A short description of your project."
		authors = [
		    {name = "Your Name", email = "your.email@example.com"}
		]
		dependencies = [
		    "package_name>=1.0",
		]

		[build-system]
		requires = ["setuptools>=61.0"]
		build-backend = "setuptools.build_meta"

		[[tool.uv.index]]
		url = "https://default.com"
		default = true

		[tool.uvbox.package]
		name = "my-package-from-pyproject"
		script = "my.script:run"
	`

	// Directory with test configuration file
	workingDirectory, _ := testDirectoryWithFileContent(t, "pyproject.toml", content)
	defer os.RemoveAll(workingDirectory)

	// Disable CLI configuration file argument
	Config = ""

	assertUvEnvironmentVariableInAutomaticallyFoundConfiguration(t, workingDirectory, "UV_DEFAULT_INDEX=https://default.com")
}

func TestAppendUvDefaultIndexNamed(t *testing.T) {
	content := `
		[project]
		name = "your_project_name"
		version = "0.1.0"
		description = "A short description of your project."
		authors = [
		    {name = "Your Name", email = "your.email@example.com"}
		]
		dependencies = [
		    "package_name>=1.0",
		]

		[build-system]
		requires = ["setuptools>=61.0"]
		build-backend = "setuptools.build_meta"

		[[tool.uv.index]]
		name = "default-index"
		url = "https://default.com"
		default = true

		[tool.uvbox.package]
		name = "my-package-from-pyproject"
		script = "my.script:run"
	`

	// Directory with test configuration file
	workingDirectory, _ := testDirectoryWithFileContent(t, "pyproject.toml", content)
	defer os.RemoveAll(workingDirectory)

	// Disable CLI configuration file argument
	Config = ""

	assertUvEnvironmentVariableInAutomaticallyFoundConfiguration(t, workingDirectory, "UV_DEFAULT_INDEX=default-index=https://default.com")
}

func TestAppendUvIndexNamed(t *testing.T) {
	content := `
		[project]
		name = "your_project_name"
		version = "0.1.0"
		description = "A short description of your project."
		authors = [
		    {name = "Your Name", email = "your.email@example.com"}
		]
		dependencies = [
		    "package_name>=1.0",
		]

		[build-system]
		requires = ["setuptools>=61.0"]
		build-backend = "setuptools.build_meta"

		[[tool.uv.index]]
		name = "a-named-index"
		url = "https://named.com"

		[tool.uvbox.package]
		name = "my-package-from-pyproject"
		script = "my.script:run"
	`

	// Directory with test configuration file
	workingDirectory, _ := testDirectoryWithFileContent(t, "pyproject.toml", content)
	defer os.RemoveAll(workingDirectory)

	// Disable CLI configuration file argument
	Config = ""

	assertUvEnvironmentVariableInAutomaticallyFoundConfiguration(t, workingDirectory, "UV_INDEX=a-named-index=https://named.com")
}

func TestAppendUvIndexUnNamed(t *testing.T) {
	content := `
		[project]
		name = "your_project_name"
		version = "0.1.0"
		description = "A short description of your project."
		authors = [
		    {name = "Your Name", email = "your.email@example.com"}
		]
		dependencies = [
		    "package_name>=1.0",
		]

		[build-system]
		requires = ["setuptools>=61.0"]
		build-backend = "setuptools.build_meta"

		[[tool.uv.index]]
		url = "https://unnamed.com"

		[tool.uvbox.package]
		name = "my-package-from-pyproject"
		script = "my.script:run"
	`

	// Directory with test configuration file
	workingDirectory, _ := testDirectoryWithFileContent(t, "pyproject.toml", content)
	defer os.RemoveAll(workingDirectory)

	// Disable CLI configuration file argument
	Config = ""

	assertUvEnvironmentVariableInAutomaticallyFoundConfiguration(t, workingDirectory, "UV_INDEX=https://unnamed.com")
}

func TestAppendMultipleIndexesWithNamedDefault(t *testing.T) {
	content := `
		[project]
		name = "your_project_name"
		version = "0.1.0"
		description = "A short description of your project."
		authors = [
		    {name = "Your Name", email = "your.email@example.com"}
		]
		dependencies = [
		    "package_name>=1.0",
		]

		[build-system]
		requires = ["setuptools>=61.0"]
		build-backend = "setuptools.build_meta"

		[[tool.uv.index]]
		name = "default-named"
		url = "https://default.com"
		default = true

		[[tool.uv.index]]
		url = "https://unnamed.com"

		[[tool.uv.index]]
		name = "named"
		url = "https://named.com"

		[tool.uvbox.package]
		name = "my-package-from-pyproject"
		script = "my.script:run"
	`

	// Directory with test configuration file
	workingDirectory, _ := testDirectoryWithFileContent(t, "pyproject.toml", content)
	defer os.RemoveAll(workingDirectory)

	// Disable CLI configuration file argument
	Config = ""

	assertUvEnvironmentVariableInAutomaticallyFoundConfiguration(t, workingDirectory, "UV_INDEX=https://unnamed.com")
	assertUvEnvironmentVariableInAutomaticallyFoundConfiguration(t, workingDirectory, "UV_INDEX=named=https://named.com")
	assertUvEnvironmentVariableInAutomaticallyFoundConfiguration(t, workingDirectory, "UV_DEFAULT_INDEX=default-named=https://default.com")
}

func TestAppendMultipleIndexesWithUnnamedDefault(t *testing.T) {
	content := `
		[project]
		name = "your_project_name"
		version = "0.1.0"
		description = "A short description of your project."
		authors = [
		    {name = "Your Name", email = "your.email@example.com"}
		]
		dependencies = [
		    "package_name>=1.0",
		]

		[build-system]
		requires = ["setuptools>=61.0"]
		build-backend = "setuptools.build_meta"

		[[tool.uv.index]]
		url = "https://default.com"
		default = true

		[[tool.uv.index]]
		url = "https://unnamed.com"

		[[tool.uv.index]]
		name = "named"
		url = "https://named.com"

		[tool.uvbox.package]
		name = "my-package-from-pyproject"
		script = "my.script:run"
	`

	// Directory with test configuration file
	workingDirectory, _ := testDirectoryWithFileContent(t, "pyproject.toml", content)
	defer os.RemoveAll(workingDirectory)

	// Disable CLI configuration file argument
	Config = ""

	assertUvEnvironmentVariableInAutomaticallyFoundConfiguration(t, workingDirectory, "UV_INDEX=https://unnamed.com")
	assertUvEnvironmentVariableInAutomaticallyFoundConfiguration(t, workingDirectory, "UV_INDEX=named=https://named.com")
	assertUvEnvironmentVariableInAutomaticallyFoundConfiguration(t, workingDirectory, "UV_DEFAULT_INDEX=https://default.com")
}

func assertUvEnvironmentVariableInAutomaticallyFoundConfiguration(t *testing.T, workingDirectory, variable string) {
	configuration, err := loadConfigurationFromEitherCLIArgumentOrDetectedFile(workingDirectory)
	if err != nil {
		t.Errorf("could not load test configuration: %s", err)
	}

	assertUvEnvironmentVariableInConfiguration(t, configuration, variable)
}

func assertUvEnvironmentVariableNotInAutomaticallyFoundConfiguration(t *testing.T, workingDirectory, variable string) {
	configuration, err := loadConfigurationFromEitherCLIArgumentOrDetectedFile(workingDirectory)
	if err != nil {
		t.Errorf("could not load test configuration: %s", err)
	}

	assertUvEnvironmentVariableNotInConfiguration(t, configuration, variable)
}

func assertUvEnvironmentVariableInConfiguration(t *testing.T, configuration Configuration, variable string) {
	if !slices.Contains(configuration.Uv.Environment, variable) {
		t.Errorf("missing '%s' from configuration uv environment variables list '%v'", variable, configuration.Uv.Environment)
	}
}

func assertUvEnvironmentVariableNotInConfiguration(t *testing.T, configuration Configuration, variable string) {
	if slices.Contains(configuration.Uv.Environment, variable) {
		t.Errorf("unexpected '%s' from configuration uv environment variables list '%v'", variable, configuration.Uv.Environment)
	}
}
