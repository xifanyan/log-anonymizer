package main

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

var (
	ListNamingPatterns = &cli.Command{
		Name:    "listNamingPatterns",
		Usage:   `log-anonymizer listNamingPatterns`,
		Aliases: []string{"ln"},
		Action:  listNamingPatterns,
	}

	ListRegexPatterns = &cli.Command{
		Name:    "listRegexPatterns",
		Usage:   `log-anonymizer listRegexPatterns`,
		Aliases: []string{"lr"},
		Action:  listRegexPatterns,
	}

	Commands = []*cli.Command{
		ListNamingPatterns,
		ListRegexPatterns,
	}
)

// listNamingPatterns lists all the naming patterns configured for the anonymizer,
// grouped by category. It calls GetNamingPatterns() on the global config to get
// the patterns, then prints them out formatted.
func listNamingPatterns(c *cli.Context) error {
	var err error

	namingPatterns, err := GlobalConfig.GetNamingPatterns(c.String("fileType"))
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

// listRegexPatterns lists all the regex patterns configured for the anonymizer,
// grouped by category. It calls GetRegexPatterns() on the global config to get
// the patterns, then prints them out formatted.
func listRegexPatterns(c *cli.Context) error {
	var err error

	regexPatterns, err := GlobalConfig.GetRegexPatterns(c.String("fileType"))
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
