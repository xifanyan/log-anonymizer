package main

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

const (
	DEFAULT_LOGTYPE     = "auto"
	DEFAULT_CONFIG      = "config.yaml"
	DEFAULT_OBFUSCATION = "_CONFIDENTIAL_"
)

func main() {
	app := &cli.App{
		Name:    "log-anonymizer",
		Version: "0.1-alpha",
		Usage:   "Axcelerate Log Anonymizer",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "logType",
				Aliases: []string{"t"},
				Usage:   "Log Type",
				Value:   DEFAULT_LOGTYPE,
			},
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "Anonymizer configuration file path",
				Value:   DEFAULT_CONFIG,
			},
			&cli.StringFlag{
				Name:    "obfuscation",
				Aliases: []string{"s"},
				Usage:   "",
				Value:   DEFAULT_OBFUSCATION,
			},
			&cli.BoolFlag{
				Name:    "debug",
				Aliases: []string{"d"},
				Usage:   "Debug Mode: Log to stderr and AdpTask log enabled",
				Value:   false,
			},
		},
		Commands: Commands,
		Before: func(c *cli.Context) error {
			zerolog.SetGlobalLevel(zerolog.ErrorLevel)
			if c.Bool("debug") {
				zerolog.SetGlobalLevel(zerolog.DebugLevel)
			}
			log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Error().Msg(err.Error())
	}
}
