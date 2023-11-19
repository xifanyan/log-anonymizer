package main

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

var (
	Version = &cli.StringFlag{
		Name:    "version",
		Aliases: []string{"v"},
		Usage:   "axcelerate version",
		Value:   "default",
	}
)

var (
	ListNamingPatterns = &cli.Command{
		Name:    "listNamingPatterns",
		Usage:   `log-anonymizer listNamingPatterns`,
		Aliases: []string{"ln"},
		Flags: []cli.Flag{
			Version,
		},
		Action: listNamingPatterns,
	}

	ListRegexes = &cli.Command{
		Name:    "listRegexes",
		Usage:   `log-anonymizer listRegexes`,
		Aliases: []string{"lr"},
		Flags: []cli.Flag{
			Version,
		},
		Action: listRegexes,
	}

	Commands = []*cli.Command{
		ListNamingPatterns,
		ListRegexes,
	}
)

func listNamingPatterns(c *cli.Context) error {
	var err error

	yaml, err := LoadConfig(c.String("config"))
	if err != nil {
		return err
	}

	cfg, err := yaml.GetAnonymizerConfigByVersion(c.String("version"))
	if err != nil {
		return err
	}

	namingPatterns, err := cfg.GetAllNamingPatterns()
	if err != nil {
		return err
	}

	i := 0
	for _, pattern := range namingPatterns {
		i++
		fmt.Printf("%-5d%-20s%s\n", i, pattern.Category, pattern.Pattern)
	}

	return err
}

func listRegexes(c *cli.Context) error {
	var err error

	yaml, err := LoadConfig(c.String("config"))
	if err != nil {
		return err
	}

	cfg, err := yaml.GetAnonymizerConfigByVersion(c.String("version"))
	if err != nil {
		return err
	}

	regexes, err := cfg.GetAllRegexes()

	i := 0
	for _, pattern := range regexes {
		i++
		fmt.Printf("%-5d%-20s%s\n", i, pattern.Category, pattern.Regex)
	}

	return err
}
