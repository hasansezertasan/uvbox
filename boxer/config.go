package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type PyProjectConfiguration struct {
	Tool ToolConfiguration `toml:"tool"`
}

type ToolConfiguration struct {
	UVBOXConfiguration Configuration       `toml:"uvbox"`
	UVConfiguration    ToolUvConfiguration `toml:"uv"`
}

type ToolUvConfiguration struct {
	Indexes             []ToolUvConfigurationIndex `toml:"index"`
	PythonInstallMirror string                     `toml:"python-install-mirror"`
}

type ToolUvConfigurationIndex struct {
	Name    string `toml:"name"`
	Url     string `toml:"url"`
	Default bool   `toml:"default"`
}

type Configuration struct {
	Package PackageConfiguration
	Certs   CertificatesConfiguration
	Uv      UvConfiguration
}

type UpdateConfiguration struct {
	Url string
}

type PackageConfiguration struct {
	Name        string
	Script      string
	Version     PackageVersionConfiguration
	Constraints PackageConstraintsConfiguration
}

type PackageVersionConfiguration struct {
	AutoUpdate bool `toml:"auto-update"`
	Static     string
	Dynamic    string
}

type PackageConstraintsConfiguration struct {
	Static  string
	Dynamic string
}

type CertificatesConfiguration struct {
	Path string
}

type UvConfiguration struct {
	Version     string
	Mirror      string
	Environment []string
}

func loadPyprojectToml(filePath string) (PyProjectConfiguration, error) {
	config := PyProjectConfiguration{}
	_, err := toml.DecodeFile(filePath, &config)
	if err != nil {
		return config, fmt.Errorf("could not decode pyproject.toml file located at %s", filePath)
	}

	return config, nil
}

func loadConfigurationFromPyprojectToml(filePath string) (Configuration, error) {
	pyprojectConfiguration, err := loadPyprojectToml(filePath)
	if err != nil {
		return Configuration{}, fmt.Errorf("could not load pyproject.toml file at %s: %w", filePath, err)
	}

	config := pyprojectConfiguration.Tool.UVBOXConfiguration

	// If possible, we determine some environment variable from uv configurations
	uvEnvironmentIndexLinesToAdd := computeConfigIndexUvEnvironmentFromToolUvIndexes(pyprojectConfiguration.Tool.UVConfiguration.Indexes)
	uvEnvironmentPythonMirrorLinesToAdd := computeConfigPythonMirrorUvEnvironmentFromUvToolIndexes(pyprojectConfiguration.Tool.UVConfiguration)
	config.Uv.Environment = append(config.Uv.Environment, uvEnvironmentIndexLinesToAdd...)
	config.Uv.Environment = append(config.Uv.Environment, uvEnvironmentPythonMirrorLinesToAdd...)

	return config, nil
}

func computeConfigPythonMirrorUvEnvironmentFromUvToolIndexes(toolUvConfig ToolUvConfiguration) (lines []string) {
	lines = []string{}

	if toolUvConfig.PythonInstallMirror != "" {
		lines = append(lines, fmt.Sprintf("UV_PYTHON_INSTALL_MIRROR=%s", toolUvConfig.PythonInstallMirror))
	}

	return lines
}

func computeConfigIndexUvEnvironmentFromToolUvIndexes(indexes []ToolUvConfigurationIndex) (lines []string) {
	lines = []string{}

	for _, index := range indexes {
		variablePart := ""
		if index.Default {
			variablePart = "UV_DEFAULT_INDEX"
		} else {
			variablePart = "UV_INDEX"
		}

		indexPart := ""
		if index.Name != "" {
			indexPart = fmt.Sprintf("=%s", index.Name)
		}

		line := fmt.Sprintf("%s%s=%s", variablePart, indexPart, index.Url)
		lines = append(lines, line)
	}

	return lines
}

func loadConfigurationFromUvboxToml(filePath string) (Configuration, error) {
	config := Configuration{}
	_, err := toml.DecodeFile(filePath, &config)
	if err != nil {
		return config, fmt.Errorf("could not decode uvbox.toml file located at %s", filePath)
	}
	return config, nil
}

func loadConfigurationFromDirectory(directory string) (Configuration, error) {
	config := Configuration{}

	// Try to find an uvbox.toml file
	if foundPath, err := directoryOrParentsContainsFile(directory, "uvbox.toml"); err != nil {
		return config, fmt.Errorf("could not check if an uvbox.toml file exists from directory '%s' or its parents: %w", directory, err)
	} else if foundPath != "" {
		return loadConfigurationFromUvboxToml(foundPath)
	}

	// Try to find a pyproject.toml file
	if foundPath, err := directoryOrParentsContainsFile(directory, "pyproject.toml"); err != nil {
		return config, fmt.Errorf("could not check if a pyproject.toml file exists from directory '%s' or its parents: %w", directory, err)
	} else if foundPath != "" {
		return loadConfigurationFromPyprojectToml(foundPath)
	}

	return config, fmt.Errorf("could not find a configuration starting from the directory '%s'", directory)
}

func loadConfigurationFromFile(filePath string) (Configuration, error) {
	name := filepath.Base(filePath)

	// We may be able to determine configuration type from the filename
	switch name {
	case "uvbox.toml":
		return loadConfigurationFromUvboxToml(filePath)
	case "pyproject.toml":
		return loadConfigurationFromPyprojectToml(filePath)
	}

	// If it's a .toml containing tool.uvbox
	if ok, err := tomlContainsToolUvbox(filePath); err != nil {
		return Configuration{}, fmt.Errorf("error while checking if the given file is an usable pyproject.toml")
	} else if ok {
		return loadConfigurationFromPyprojectToml(filePath)
	} else {
		return loadConfigurationFromUvboxToml(filePath)
	}
}

func loadConfigurationFromEitherCLIArgumentOrDetectedFile(workingDirectory string) (Configuration, error) {
	if Config != "" {
		return loadConfigurationFromFile(Config)
	} else {
		return loadConfigurationFromDirectory(workingDirectory)
	}
}

func loadConfigurationFromCLIOrCurrentDirectory() (Configuration, error) {
	// Get current directory
	workingDirectory, err := os.Getwd()
	if err != nil {
		return Configuration{}, fmt.Errorf("could not get working directory: %w", err)
	}

	return loadConfigurationFromEitherCLIArgumentOrDetectedFile(workingDirectory)
}

func saveConfigurationToDirectory(configuration Configuration, directory, fileName string) error {
	// Serialize configuration
	tomlData, err := toml.Marshal(configuration)
	if err != nil {
		return fmt.Errorf("could not marshal configuration: %w", err)
	}

	// Write file
	targetPath := filepath.Join(directory, fileName)
	if err := os.WriteFile(targetPath, tomlData, 0644); err != nil {
		return fmt.Errorf("could not write configuration file into %s: %w", targetPath, err)
	}

	return nil
}

func (c Configuration) ParsedPackageConfiguration() (PackageConfiguration, error) {
	// Validate name
	if c.Package.Name == "" {
		return c.Package, fmt.Errorf("package.name must not be empty")
	}

	// Validate script
	if c.Package.Script == "" {
		return c.Package, fmt.Errorf("package.script must not be empty")
	}

	return c.Package, nil
}

func configurationScriptName() (string, error) {
	// Load configuration file
	config, err := loadConfigurationFromCLIOrCurrentDirectory()
	if err != nil {
		return "", fmt.Errorf("could not load configuration: %w", err)
	}

	// Load package configuration with validated values
	if packageConfiguration, err := config.ParsedPackageConfiguration(); err != nil {
		return "", fmt.Errorf("could not load configuration package script name: %w", err)
	} else {
		return packageConfiguration.Script, nil
	}
}

func configurationCertificatesBundlePath() (string, error) {
	// Load configuration file
	config, err := loadConfigurationFromCLIOrCurrentDirectory()
	if err != nil {
		return "", fmt.Errorf("could not load configuration: %w", err)
	}

	return config.Certs.Path, nil
}

func directoryOrParentsContainsFile(directory string, fileName string) (string, error) {
	// Will use absolute paths of given directory to start and to be able to correctly iterate
	startDirectory, err := filepath.Abs(directory)
	if err != nil {
		return "", fmt.Errorf("could get absolute path of directory %s: %w", directory, err)
	}

	// Look for given file in directory, and recursively update directory to parent one if not found
	for dir := startDirectory; dir != ""; dir = filepath.Dir(dir) {
		if foundPath, err := directoryContainsFile(dir, fileName); err != nil {
			return foundPath, fmt.Errorf("could not check if file exists: %w", err)
		} else if foundPath != "" {
			return foundPath, nil
		} else if dir == filepath.Dir(dir) { // We are at the root of the filesystem
			break
		}
	}

	return "", nil
}

func directoryContainsFile(directory string, fileName string) (string, error) {
	filePath := filepath.Join(directory, fileName)
	if ok, err := fileExists(filePath); err != nil {
		return "", fmt.Errorf("error while checking if directory '%s' contains file '%s': %w", directory, fileName, err)
	} else if ok {
		return filePath, nil
	} else {
		return "", nil
	}
}

func fileExists(fileName string) (bool, error) {
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("could not check if file exists at %s: %w", fileName, err)
	} else {
		return true, nil
	}
}

func tomlContainsToolUvbox(tomlFile string) (bool, error) {
	config := PyProjectConfiguration{}
	md, err := toml.DecodeFile(tomlFile, &config)
	if err != nil {
		return false, fmt.Errorf("could not decode pyproject.toml file located at %s: %w", tomlFile, err)
	}

	return md.IsDefined("tool", "uvbox"), nil
}
