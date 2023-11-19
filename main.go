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
	DEFAULT_AXC_VERSION = "default"
	DEFAULT_OBFUSCATION = "_CONFIDENTIAL_"
)

var GlobalConfig *AnonymizerConfig

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
				Name:    "axcVersion",
				Aliases: []string{"a"},
				Usage:   "axcelerate version",
				Value:   DEFAULT_AXC_VERSION,
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
			var err error

			zerolog.SetGlobalLevel(zerolog.ErrorLevel)
			if c.Bool("debug") {
				zerolog.SetGlobalLevel(zerolog.DebugLevel)
			}
			log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

			log.Debug().Msgf("axcVersion: %+v", c.String("axcVersion"))

			yamlCfg, err := LoadConfig(c.String("config"))
			if err != nil {
				return err
			}
			log.Debug().Msgf("yaml Config: %+v", yamlCfg)

			// DO NOT USE := to set global variable due to variable shawdowing
			GlobalConfig, err = yamlCfg.GetAnonymizerConfigByAxcVersion(c.String("axcVersion"))
			if err != nil {
				return err
			}
			log.Debug().Msgf("GlobalConfig: %+v", GlobalConfig)

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Error().Msg(err.Error())
	}
}
