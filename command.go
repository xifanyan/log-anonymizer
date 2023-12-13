package main

import (
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
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

	CleanUp = &cli.Command{
		Name:    "cleanUp",
		Usage:   `log-anonymizer cleanUp`,
		Aliases: []string{"cu"},
		Flags: []cli.Flag{
			Path,
		},
		Action: cleanUp,
	}

	Run = &cli.Command{
		Name:  "run",
		Usage: `log-anonymizer run --path ./service.log`,
		Flags: []cli.Flag{
			Path,
			WorkerCount,
		},
		Action: run,
	}

	Commands = []*cli.Command{
		CleanUp,
		ListNamingPatterns,
		ListRegexPatterns,
		ListKinds,
		Run,
	}
)

var (
	Path = &cli.StringFlag{
		Name:     "path",
		Usage:    "file or folder to be processed",
		Required: true,
	}

	WorkerCount = &cli.IntFlag{
		Name:  "workerCount",
		Usage: "number of workers",
		Value: 2,
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

/**
* Run - Process files or folders based on command-line flags.
* Inputs:
*	c: A cli.Context object that contains the command-line context and flags.
*
* Outputs:
*	err (error): A error that occurred during process.
 */
func run(c *cli.Context) error {
	scheduler := NewScheduler().
		WithPath(c.String("path")).
		WithKind(c.String("kind")). // defined as global flag in main.go
		WithWorkerCount(c.Int("workerCount")).
		WithObfuscation(c.String("obfuscation"))

	filePaths, err := scheduler.getLogs()
	if err != nil {
		return err
	}

	scheduler.Process(filePaths)

	return nil
}

// cleanUp deletes the anonymized log files after processing is complete.
// It gets the slice of anonymized file paths from the context
// and calls deleteFiles to delete each one.
//
// Parameters:
//   - c: The CLI context containing the anonymized files slice
//
// Returns:
//   - error: Any error encountered while deleting files
func cleanUp(c *cli.Context) error {
	var err error

	scheduler := NewScheduler().
		WithPath(c.String("path"))

	anonymizedLogs, err := scheduler.getAnonymizedLogs()
	if err != nil {
		return err
	}

	if len(anonymizedLogs) == 0 {
		log.Info().Msg("No anonymized logs found")
		return nil
	}

	for _, anonymized := range anonymizedLogs {
		log.Debug().Msgf("%s", anonymized)
		err = os.Remove(anonymized)
		if err != nil {
			log.Error().Msgf("%s", err)
			continue
		}
	}
	return nil
}
