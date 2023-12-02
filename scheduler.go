package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sync"

	"github.com/rs/zerolog/log"
)

type Scheduler struct {
	path        string
	kind        string
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

type logFileInfo struct {
	path string
	kind string
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
				absPath, err := filepath.Abs(path)
				if err != nil {
					log.Error().Msgf("%s", err)
				} else {
					infos = append(infos, logFileInfo{path: absPath, kind: kind})
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

func (s *Scheduler) processFile(info logFileInfo) error {
	fmt.Println("Processing file:", info.path, "of kind:", info.kind)
	return nil
}
