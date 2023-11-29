package main

import (
	"fmt"
	"os"
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

type logInfo struct {
	path string
	kind string
}

func (s *Scheduler) getKindByPath() (string, error) {
	return "", nil
}

// getLogs retrieves log information from the specified path.
//
// This function walks through the directory tree rooted at the 'path' field of the Scheduler struct.
// For each file in the directory tree, it checks if the file is not a directory.
// If the 'kind' field of the Scheduler struct is "*", it identifies the kind of log type.
// It then retrieves the absolute path of the file and appends it along with the identified log type to the 'logInfos' slice.
//
// Returns:
//   - []logInfo: a slice of logInfo structs containing the path and kind of each log file.
//   - error: an error, if one occurred during the filepath.Walk function or during the retrieval of the absolute path.

func (s *Scheduler) getLogs() ([]logInfo, error) {
	var logInfos []logInfo

	err := filepath.Walk(s.path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			var kind string = "NONE"
			if s.kind == "*" {
				// kind of log type is not specified which requires to be identified
				kind, _ = s.getKindByPath()
			}

			absPath, err := filepath.Abs(path)
			if err != nil {
				return err
			}
			logInfos = append(logInfos, logInfo{path: absPath, kind: kind})
		}
		return nil
	})

	return logInfos, err
}

// Process processes log files using a worker pool.
//
// This function creates a worker pool of goroutines to process log files.
// Each worker goroutine takes a logInfo from the 'pathChan' channel, processes it using the 'processFile' method,
// and logs any errors that occur during processing.
//
// Parameters:
//   - logInfos: a slice of logInfo structs containing the path and kind of each log file.
//
// This function does not return any values.
func (s *Scheduler) Process(logInfos []logInfo) {
	var wg sync.WaitGroup
	pathChan := make(chan logInfo, s.workerCount)

	for i := 0; i < s.workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for fi := range pathChan {
				if err := s.processFile(fi); err != nil {
					log.Error().Msgf("%s", err)
				}
			}
		}()
	}

	for _, logInfo := range logInfos {
		pathChan <- logInfo
	}
	close(pathChan)

	wg.Wait()
}

func (s *Scheduler) processFile(logInfo logInfo) error {
	fmt.Println("Processing file:", logInfo.path)
	return nil
}
