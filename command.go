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

	Commands = []*cli.Command{
		ListNamingPatterns,
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

	for _, pattern := range namingPatterns {
		fmt.Printf("%-20s%s\n", pattern.Category, pattern.Pattern)
	}

	return err
}
