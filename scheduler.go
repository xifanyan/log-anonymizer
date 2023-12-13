package main

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

type Scheduler struct {
	path        string
	kind        string
	obfuscation string
	workerCount int
}

func NewScheduler() *Scheduler {
	return &Scheduler{}
}

// WithPath sets the path field of the Scheduler object to the provided path parameter and returns a pointer to the modified Scheduler object.
//
// Parameters:
// - path (string): The path to be set for the Scheduler object.
//
// Returns:
// - *Scheduler: A pointer to the modified Scheduler object.
func (s *Scheduler) WithPath(path string) *Scheduler {
	s.path = path
	return s
}

// WithKind sets the kind field of the Scheduler to the provided kind parameter and returns a pointer to the modified Scheduler.
//
// Parameters:
// - kind (string): The kind to be set for the Scheduler object.
//
// Returns:
// - *Scheduler: A pointer to the modified Scheduler object.

func (s *Scheduler) WithKind(kind string) *Scheduler {
	s.kind = kind
	return s
}

// WithWorkerCount sets the workerCount field of the Scheduler object to the provided workerCount parameter and returns a pointer to the modified Scheduler object.
//
// Parameters:
// - workerCount (int): The number of workers to be set for the Scheduler object.
//
// Returns:
// - *Scheduler: A pointer to the modified Scheduler object.
func (s *Scheduler) WithWorkerCount(workerCount int) *Scheduler {
	s.workerCount = workerCount
	return s
}

// WithObfuscation sets the obfuscation field of the Scheduler object to the provided obfuscation parameter and returns a pointer to the modified Scheduler object.
//
// Parameters:
// - obfuscation (string): the obfuscation string.
//
// Returns:
// - *Scheduler: A pointer to the modified Scheduler object.
func (s *Scheduler) WithObfuscation(obfuscation string) *Scheduler {
	s.obfuscation = obfuscation
	return s
}

type logFileInfo struct {
	kind string
	path string
}

// getOutputFileName returns the output file name for the log file info.
// It appends ".anonymized" and timestamp to the end of the log file path.
//
// Returns:
//   - string: The output file name for the log info.
func (inf logFileInfo) getOutputFileName() string {
	now := time.Now()
	ts := fmt.Sprintf("%d%02d%02d-%02d%02d%02d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
	outputFileName := fmt.Sprintf("%s.anonymized.%s", inf.path, ts)
	return outputFileName
}

// getKindByLogPath identifies the kind of log file based on its path and configured naming patterns.
// It takes the scheduler path, gets naming patterns for the scheduler kind, and checks if the base log file name matches any naming pattern.
//
// Parameters:
// - logFilePath (string): log file path
//
// Returns:
//   - string: The type of log file identified by the provided path.
//   - error: An error indicating if there were any issues retrieving the type of log file.
func (s *Scheduler) getKindByLogPath(logFilePath string) (string, error) {
	var err error

	namingPatterns, err := GlobalConfig.GetNamingPatterns(s.kind)
	if err != nil {
		return "", err
	}

	// get base log file name
	logName := path.Base(logFilePath)

	for _, namingPattern := range namingPatterns {
		if namingPattern.Regex.MatchString(logName) {
			return namingPattern.Kind, nil
		}
	}

	err = fmt.Errorf("not able to detect log type: %s", logName)
	return "", err
}

// getLogs retrieves log information from the specified path.
// It walks the provided path and collects information about each log file found.
// For each log file, it determines the kind of log based on configured naming patterns.
//
// Returns:
//   - []logInfo: a slice of logInfo structs containing the path and kind of each log file.
//   - error: an error, if one occurred during the filepath.Walk function or during the retrieval of the absolute path.
func (s *Scheduler) getLogs() ([]logFileInfo, error) {
	var infos []logFileInfo

	err := filepath.Walk(s.path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			var kind string

			if s.kind == "*" {
				kind, err = s.getKindByLogPath(path)
			} else {
				kind = s.kind
			}

			if err != nil {
				log.Error().Msgf("%s", err)
			} else {
				if !strings.Contains(path, ".anonymized.") {
					absPath, err := filepath.Abs(path)
					if err != nil {
						log.Error().Msgf("%s", err)
					} else {
						infos = append(infos, logFileInfo{path: absPath, kind: kind})
					}
				}
			}

		}
		return nil
	})

	return infos, err
}

// Process processes log files using a worker pool.
//
// This function creates a worker pool of goroutines to process log files.
// Each worker goroutine takes a logInfo from the 'pathChan' channel, processes it using the 'processFile' method,
// and logs any errors that occur during processing.
//
// Parameters:
//   - infos ([]LogFileInfo): a slice of logFileInfo structs containing the path and kind of each log file.
//
// This function does not return any values.
func (s *Scheduler) Process(infos []logFileInfo) {
	var wg sync.WaitGroup
	pathChan := make(chan logFileInfo, s.workerCount)

	for i := 0; i < s.workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for info := range pathChan {
				if err := s.processFile(info); err != nil {
					log.Error().Msgf("%s", err)
				}
			}
		}()
	}

	for _, info := range infos {
		pathChan <- info
	}
	close(pathChan)

	wg.Wait()
}

// processFile processes an individual log file.
// It reads the log file line by line, applies the configured obfuscation,
// and writes the anonymized output to a new file.
//
// Parameters:
//   - info: logFileInfo containing log file path and type
//
// Returns:
//   - error: any error encountered while processing the file
func (s *Scheduler) processFile(info logFileInfo) error {
	var err error

	log.Debug().Msgf("processing [%s] log file: %s", info.kind, info.path)

	regexes, err := GlobalConfig.GetRegexPatterns(info.kind)
	if err != nil {
		return err
	}

	inf, err := os.Open(info.path)
	if err != nil {
		return err
	}
	defer inf.Close()

	outf, err := os.Create(info.getOutputFileName())
	if err != nil {
		return err
	}
	defer outf.Close()

	fs := bufio.NewScanner(inf)

	wf := bufio.NewWriter(outf)

	for fs.Scan() {
		line := s.obfuscate(fs.Text(), regexes)

		_, err := fmt.Fprintln(wf, line)
		if err != nil {
			log.Error().Msgf("%s", err)
			return err
		}
	}

	if err := fs.Err(); err != nil {
		log.Error().Msgf("%s", err)
		return err
	}

	log.Debug().Msgf("finished processing [%s] log file: %s", info.kind, info.path)

	return nil
}

// obfuscate takes a log line and a slice of obfuscation patterns, and returns
// the log line with sensitive information obfuscated.
//
// It iterates through each obfuscation pattern and replaces any matches in
// the log line with the pattern's obfuscation string.
//
// Parameters:
//   - line: the log line to obfuscate
//   - regexes: slice of obfuscation patterns to apply
//
// Returns:
//   - The obfuscated log line
func (s *Scheduler) obfuscate(line string, regexes []Pattern) string {
	for _, re := range regexes {
		line = re.Regex.ReplaceAllStringFunc(line, func(matched string) string {
			matches := re.Regex.FindStringSubmatch(matched)
			log.Debug().Msgf("- [matched] %s", line)
			for _, match := range matches[1:] {
				matched = strings.Replace(matched, match, s.obfuscation, 1)
			}
			return matched
		})
	}
	return line
}

// getAnonymizedLogs collect anonymized log file names.
//
// Returns:
//   - []string: Slice of paths for the anonymized log files
//   - error: Any error encountered
func (s *Scheduler) getAnonymizedLogs() ([]string, error) {

	var paths []string

	err := filepath.Walk(s.path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			if strings.Contains(path, ".anonymized.") {
				absPath, _ := filepath.Abs(path)
				paths = append(paths, absPath)
			}
		}
		return nil
	})

	return paths, err
}
