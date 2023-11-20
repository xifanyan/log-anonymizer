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

	ListKinds = &cli.Command{
		Name:    "listKinds",
		Usage:   `log-anonymizer listKinds`,
		Aliases: []string{"lk"},
		Action:  listKinds,
	}

	Commands = []*cli.Command{
		ListNamingPatterns,
		ListRegexPatterns,
		ListKinds,
	}
)

// listNamingPatterns lists all the naming patterns configured for the anonymizer.
// It prints out the list of patterns with index numbers. The patterns can be
// filtered by kind using the --kind flag.
func listNamingPatterns(c *cli.Context) error {
	var err error

	namingPatterns, err := GlobalConfig.GetNamingPatterns(c.String("kind"))
	if err != nil {
		return err
	}

	for i, pattern := range namingPatterns {
		fmt.Printf("%-4d%-16s%s\n", i+1, pattern.Kind, pattern.Pattern)
	}

	return err
}

// listRegexPatterns lists all the regex patterns configured for the anonymizer.
// It prints out the list of regexes with index numbers. The regexes can be
// filtered by kind using the --kind flag.
func listRegexPatterns(c *cli.Context) error {
	var err error

	regexPatterns, err := GlobalConfig.GetRegexPatterns(c.String("kind"))
	if err != nil {
		return err
	}

	for i, pattern := range regexPatterns {
		fmt.Printf("%-4d%-16s%s\n", i+1, pattern.Kind, pattern.Regex)
	}

	return err
}

// listKinds lists all the log types configured for the anonymizer.
// It prints out the list of kinds with index numbers.
func listKinds(c *cli.Context) error {
	var err error

	kinds, err := GlobalConfig.GetKinds()
	if err != nil {
		return err
	}

	for i, kind := range kinds {
		fmt.Printf("%-4d%s\n", i+1, kind)
	}

	return err
}
