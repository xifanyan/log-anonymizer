package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type AnonymizerConfiguration struct {
	Configs []AnonymizerConfig `yaml:"anonymizer"`
}

type AnonymizerConfig struct {
	Version    string      `yaml:"version"` // Version identifier
	LogConfigs []LogConfig `yaml:"logs"`
}

type LogConfig struct {
	Category       string   `yaml:"category"`
	NamingPatterns []string `yaml:"namingPattern"`
	Regexes        []string `yaml:"regexes"`
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
func (cfg *AnonymizerConfiguration) GetAnonymizerConfigByVersion(version string) (*AnonymizerConfig, error) {
	for _, log := range cfg.Configs {
		if log.Version == version {
			return &log, nil
		}
	}
	return nil, fmt.Errorf("no config found for version %s", version)
}

type NamingPattern struct {
	Category string
	Pattern  string
}

// GetAllNamingPatterns returns all the naming patterns configured for the anonymizer,
// grouped by category. It loops through all the log configs and extracts the naming
// patterns into a slice of NamingPattern structs.
func (cfg *AnonymizerConfig) GetAllNamingPatterns() ([]NamingPattern, error) {
	var namingPatterns = []NamingPattern{}

	for _, log := range cfg.LogConfigs {
		for _, pattern := range log.NamingPatterns {
			namingPatterns = append(namingPatterns, NamingPattern{
				Category: log.Category,
				Pattern:  pattern,
			})
		}
	}

	if len(namingPatterns) == 0 {
		return namingPatterns, fmt.Errorf("no naming patterns found for %s", cfg.Version)
	}

	return namingPatterns, nil
}

// GetLogConfigByCategory returns the LogConfig for the given category from the
// AnonymizerConfig. It loops through all the LogConfigs and returns the one
// where the Category matches the passed in category string. If no match is found,
// it returns an error.
func (cfg *AnonymizerConfig) GetLogConfigByCategory(category string) (*LogConfig, error) {
	for _, logCfg := range cfg.LogConfigs {
		if logCfg.Category == category {
			return &logCfg, nil
		}
	}
	return nil, fmt.Errorf("no config found for category %s under %s", category, cfg.Version)
}
