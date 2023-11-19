package main

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

var (
	ListAllNamingPatterns = &cli.Command{
		Name:    "listAllNamingPatterns",
		Usage:   `log-anonymizer listAllNamingPatterns`,
		Aliases: []string{"ln"},
		Action:  listAllNamingPatterns,
	}

	ListAllRegexPatterns = &cli.Command{
		Name:    "listAllRegexPatterns",
		Usage:   `log-anonymizer listAllRegexPatterns`,
		Aliases: []string{"lr"},
		Action:  listAllRegexPatterns,
	}

	Commands = []*cli.Command{
		ListAllNamingPatterns,
		ListAllRegexPatterns,
	}
)

// listAllNamingPatterns lists all the naming patterns configured for the anonymizer,
// grouped by category. It calls GetAllNamingPatterns() on the global config to get
// the patterns, then prints them out formatted.
func listAllNamingPatterns(c *cli.Context) error {
	var err error

	namingPatterns, err := GlobalConfig.GetAllNamingPatterns()
	if err != nil {
		return err
	}

	i := 0
	for _, pattern := range namingPatterns {
		i++
		fmt.Printf("%-3d%-16s%s\n", i, pattern.Category, pattern.Pattern)
	}

	return err
}

// listAllRegexPatterns lists all the regex patterns configured for the anonymizer,
// grouped by category. It calls GetAllRegexPatterns() on the global config to get
// the patterns, then prints them out formatted.
func listAllRegexPatterns(c *cli.Context) error {
	var err error

	regexPatterns, err := GlobalConfig.GetAllRegexPatterns()
	if err != nil {
		return err
	}

	i := 0
	for _, pattern := range regexPatterns {
		i++
		fmt.Printf("%-3d%-16s%s\n", i, pattern.Category, pattern.Regex)
	}

	return err
}
