package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type AnonymizerConfiguration struct {
	AnonymizerConfigs []AnonymizerConfig `yaml:"anonymizer"`
}

type AnonymizerConfig struct {
	AxcVersion string      `yaml:"axcVersion"` // axcelerate version
	LogConfigs []LogConfig `yaml:"logs"`
}

type LogConfig struct {
	FileType       string   `yaml:"fileType"`
	NamingPatterns []string `yaml:"namingPatterns"`
	RegexPatterns  []string `yaml:"regexPatterns"`
}

// LoadConfig loads the anonymizer configuration from the YAML file at the given path.
// It returns a pointer to an AnonymizerConfiguration struct containing the parsed config,
// and an error if there were any issues reading or parsing the file.
func LoadConfig(path string) (*AnonymizerConfiguration, error) {
	var err error
	var cfg = new(AnonymizerConfiguration)

	yamlFile, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Error reading YAML file:", err)
		return cfg, err
	}

	err = yaml.Unmarshal(yamlFile, cfg)
	if err != nil {
		fmt.Println("Error unmarshaling YAML:", err)
		return cfg, err
	}

	return cfg, nil
}

// GetAnonymizerConfigByVersion returns the AnonymizerConfig for the given version from the
// AnonymizerConfiguration. It loops through all the Configs and returns the one
// where the Version matches the passed in version string. If no match is found,
// it returns an error.
func (cfg *AnonymizerConfiguration) GetAnonymizerConfigByAxcVersion(version string) (*AnonymizerConfig, error) {
	for _, anonymizerCfg := range cfg.AnonymizerConfigs {
		if anonymizerCfg.AxcVersion == version {
			return &anonymizerCfg, nil
		}
	}
	return nil, fmt.Errorf("no config found for version %s", version)
}

type NamingPattern struct {
	Category string
	Pattern  string
}

// GetNamingPatterns returns all the naming patterns configured for the anonymizer
// that match the given category. It loops through all the LogConfigs and extracts
// the naming patterns into a slice of NamingPattern structs. Pass "*" to get all
// naming patterns across all categories.
func (cfg *AnonymizerConfig) GetNamingPatterns(fileType string) ([]NamingPattern, error) {
	var namingPatterns = []NamingPattern{}

	for _, logCfg := range cfg.LogConfigs {
		for _, pattern := range logCfg.NamingPatterns {
			if fileType == logCfg.FileType || fileType == "*" {
				namingPatterns = append(namingPatterns, NamingPattern{
					Category: logCfg.FileType,
					Pattern:  pattern,
				})
			}
		}
	}

	if len(namingPatterns) == 0 {
		return namingPatterns, fmt.Errorf("no naming patterns found for %s", cfg.AxcVersion)
	}

	return namingPatterns, nil
}

type RegexPattern struct {
	Category string
	Regex    string
}

// GetRegexPatterns returns all the regex patterns configured for the anonymizer
// that match the given category. It loops through all the LogConfigs and extracts
// the regex patterns into a slice of RegexPattern structs. Pass "*" to get all
// regex patterns across all categories.
func (cfg *AnonymizerConfig) GetRegexPatterns(fileType string) ([]RegexPattern, error) {
	var regexPatterns = []RegexPattern{}

	for _, logCfg := range cfg.LogConfigs {
		for _, pattern := range logCfg.RegexPatterns {
			if fileType == logCfg.FileType || fileType == "*" {
				regexPatterns = append(regexPatterns, RegexPattern{
					Category: logCfg.FileType,
					Regex:    pattern,
				})
			}
		}
	}

	if len(regexPatterns) == 0 {
		return regexPatterns, fmt.Errorf("no regexes found for %s", cfg.AxcVersion)
	}

	return regexPatterns, nil
}

// GetLogConfigByFileType returns the LogConfig for the given fileType.
// It loops through all the configs and returns the one where FileType matches.
// Returns error if no matching config is found.
func (cfg *AnonymizerConfig) GetLogConfigByFileType(fileType string) (*LogConfig, error) {
	for _, logCfg := range cfg.LogConfigs {
		if logCfg.FileType == fileType {
			return &logCfg, nil
		}
	}
	return nil, fmt.Errorf("no config found for file type %s under %s", fileType, cfg.AxcVersion)
}
