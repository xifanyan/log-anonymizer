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

/**
 * listNamingPatterns is a function that takes a cli.Context as input and returns an error.
 * It retrieves naming patterns from a global configuration based on the value of the "kind" flag provided in the context.
 * It then iterates over the retrieved patterns and prints them out with their corresponding index, kind, and pattern.
 *
 * Inputs:
 *     - c: A cli.Context object that contains the command-line context and flags.
 *
 * Outputs:
 *     - Prints out the index, kind, and pattern of each naming pattern.
 */
func listNamingPatterns(c *cli.Context) error {
	namingPatterns, err := GlobalConfig.GetNamingPatterns(c.String("kind"))
	if err != nil {
		return err
	}

	for i, pattern := range namingPatterns {
		fmt.Printf("%-4d%-16s%s\n", i+1, pattern.Kind, pattern.Pattern)
	}

	return nil
}

/**
 * listRegexPatterns is a function that takes a cli.Context object as input and returns an error.
 * It retrieves regex patterns from a global configuration based on the value of the "kind" flag provided in the context.
 * It then iterates over the retrieved patterns and prints them out with their corresponding index, kind, and pattern.
 *
 * Inputs:
 *     - c: A cli.Context object that contains the command-line context and flags.
 *
 * Outputs:
 *     - Prints out the index, kind, and pattern of each regex pattern.
 */
func listRegexPatterns(c *cli.Context) error {
	regexPatterns, err := GlobalConfig.GetRegexPatterns(c.String("kind"))
	if err != nil {
		return err
	}

	for i, pattern := range regexPatterns {
		fmt.Printf("%-4d%-16s%s\n", i+1, pattern.Kind, pattern.Pattern)
	}

	return nil
}

/**
 * listKinds is a function that takes a cli.Context object as input and returns an error. It retrieves kinds from a global configuration and prints them out with their corresponding index.
 *
 * Inputs:
 *     - c: A cli.Context object that contains the command-line context and flags.
 *
 * Outputs:
 *     - Prints out the index and kind of each kind.
 */
func listKinds(c *cli.Context) error {
	kinds, err := GlobalConfig.GetKinds()
	if err != nil {
		return err
	}

	for i, kind := range kinds {
		fmt.Printf("%-4d%s\n", i+1, kind)
	}

	return nil
}
