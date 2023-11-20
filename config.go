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
	Kind           string   `yaml:"kind"`
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
			if fileType == logCfg.Kind || fileType == "*" {
				namingPatterns = append(namingPatterns, NamingPattern{
					Category: logCfg.Kind,
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
func (cfg *AnonymizerConfig) GetRegexPatterns(kind string) ([]RegexPattern, error) {
	var regexPatterns = []RegexPattern{}

	for _, logCfg := range cfg.LogConfigs {
		for _, pattern := range logCfg.RegexPatterns {
			if kind == logCfg.Kind || kind == "*" {
				regexPatterns = append(regexPatterns, RegexPattern{
					Category: logCfg.Kind,
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

// GetKinds returns all the log kinds configured for the anonymizer.
// It loops through the LogConfigs slice in the AnonymizerConfig and extracts
// the Kind into a string slice.
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

// GetLogConfigByLogType returns the LogConfig struct for the given logType.
// It loops through the LogConfigs slice in the AnonymizerConfig and returns
// the one where LogType matches the passed in logType. If no match is found,
// it returns an error.
func (cfg *AnonymizerConfig) GetLogConfigByLogType(kind string) (*LogConfig, error) {
	for _, logCfg := range cfg.LogConfigs {
		if logCfg.Kind == kind {
			return &logCfg, nil
		}
	}
	return nil, fmt.Errorf("no config found for file type %s under %s", kind, cfg.AxcVersion)
}
