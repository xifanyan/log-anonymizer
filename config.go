package main

import (
	"fmt"
	"os"
	"regexp"

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
	Kind           string   `yaml:"kind"`
	NamingPatterns []string `yaml:"namingPatterns"`
	RegexPatterns  []string `yaml:"regexPatterns"`
}

// LoadConfig loads a YAML configuration file and unmarshals it into an AnonymizerConfiguration struct.
//
// Inputs:
// - path (string): The path to the YAML configuration file.
//
// Outputs:
// - *AnonymizerConfiguration: A pointer to the loaded and unmarshaled AnonymizerConfiguration struct.
// - error: An error indicating if there were any issues loading or unmarshaling the YAML configuration file.
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

// GetAnonymizerConfigByAxcVersion searches for an AnonymizerConfig object in the AnonymizerConfigs slice based on the provided version parameter.
//
// Inputs:
// - version (string): The version of the AnonymizerConfig to retrieve.
//
// Outputs:
// - *AnonymizerConfig: A pointer to the matching AnonymizerConfig object if found, or nil if not found.
// - error: An error indicating if there were any issues retrieving the AnonymizerConfig.
func (cfg *AnonymizerConfiguration) GetAnonymizerConfigByAxcVersion(version string) (*AnonymizerConfig, error) {
	for _, anonymizerCfg := range cfg.AnonymizerConfigs {
		if anonymizerCfg.AxcVersion == version {
			return &anonymizerCfg, nil
		}
	}
	return nil, fmt.Errorf("no config found for version %s", version)
}

type Pattern struct {
	Kind    string
	Pattern string
	Regex   *regexp.Regexp
}

// GetNamingPatterns retrieves a list of naming patterns based on the provided kind parameter.
//
// Inputs:
//   - kind (string): The kind of log for which to retrieve the naming patterns.
//
// Outputs:
//   - []Pattern: A slice of Pattern structs containing the retrieved naming patterns.
//   - error: An error indicating if there were any issues retrieving the naming patterns.
func (cfg *AnonymizerConfig) GetNamingPatterns(kind string) ([]Pattern, error) {
	var namingPatterns = []Pattern{}

	for _, logCfg := range cfg.LogConfigs {
		for _, pattern := range logCfg.NamingPatterns {
			rex, _ := regexp.Compile(pattern)
			if kind == logCfg.Kind || kind == "*" {
				namingPatterns = append(namingPatterns, Pattern{
					Kind:    logCfg.Kind,
					Pattern: pattern,
					Regex:   rex,
				})
			}
		}
	}

	if len(namingPatterns) == 0 {
		return namingPatterns, fmt.Errorf("no naming patterns found for %s", cfg.AxcVersion)
	}

	return namingPatterns, nil
}

// GetRegexPatterns retrieves a list of regex patterns based on the provided kind parameter.
//
// Inputs:
//   - kind (string): The kind of log for which to retrieve the regex patterns.
//
// Outputs:
//   - []Pattern: A slice of Pattern structs containing the retrieved regex patterns.
//   - error: An error indicating if there were any issues retrieving the regex patterns.
func (cfg *AnonymizerConfig) GetRegexPatterns(kind string) ([]Pattern, error) {
	var regexPatterns = []Pattern{}

	for _, logCfg := range cfg.LogConfigs {
		for _, pattern := range logCfg.RegexPatterns {
			if kind == logCfg.Kind || kind == "*" {
				rex, _ := regexp.Compile(pattern)
				regexPatterns = append(regexPatterns, Pattern{
					Kind:    logCfg.Kind,
					Pattern: pattern,
					Regex:   rex,
				})
			}
		}
	}

	if len(regexPatterns) == 0 {
		return regexPatterns, fmt.Errorf("no regexes found for %s", cfg.AxcVersion)
	}

	return regexPatterns, nil
}

// GetKinds retrieves a list of log types (kinds) from the LogConfigs slice in the AnonymizerConfig struct.
//
// Inputs:
//   - None
//
// Outputs:
//   - kinds: A slice of strings containing the log types (kinds) found in the LogConfigs slice of the AnonymizerConfig.
//   - err: An error indicating if no log types were found.
func (cfg *AnonymizerConfig) GetKinds() ([]string, error) {
	var err error

	var kinds = []string{}
	for _, logCfg := range cfg.LogConfigs {
		kinds = append(kinds, logCfg.Kind)
	}

	if len(kinds) == 0 {
		return kinds, fmt.Errorf("no log types found for %s", cfg.AxcVersion)
	}

	return kinds, err
}

// GetLogConfigByLogType searches for a LogConfig with a matching Kind value in the LogConfigs slice of the AnonymizerConfig.
// Inputs:
//   - kind (string): The kind of log for which to retrieve the regex patterns.
//
// Outputs:
//   - *LogConfig: pointer to the matching LogConfig
//   - error: An error indicating if there were any issues retrieving the LogConfig
func (cfg *AnonymizerConfig) GetLogConfigByLogType(kind string) (*LogConfig, error) {
	for _, logCfg := range cfg.LogConfigs {
		if logCfg.Kind == kind {
			return &logCfg, nil
		}
	}
	return nil, fmt.Errorf("no config found for file type %s under %s", kind, cfg.AxcVersion)
}
