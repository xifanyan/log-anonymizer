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

// Traverse traverses a directory and returns a list of file paths.
// Outputs:
//
//	filePaths ([]string): A slice containing the absolute paths of all the files in the directory (including subdirectories).
//	err (error): An error that occurred during the traversal, if any.
func (s *Scheduler) Traverse() ([]string, error) {
	var filePaths []string

	err := filepath.Walk(s.path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			absPath, err := filepath.Abs(path)
			if err != nil {
				return err
			}
			filePaths = append(filePaths, absPath)
		}
		return nil
	})

	return filePaths, err
}

// Process is a method belonging to the Scheduler struct. It takes a slice of file paths as input and processes each file concurrently using goroutines and channels.
//
// Inputs:
// - filePaths (slice of strings): A list of file paths to be processed.
//
// Outputs:
// None. The method processes the files concurrently but does not return any output.
func (s *Scheduler) Process(filePaths []string) {
	var wg sync.WaitGroup
	pathChan := make(chan string, s.workerCount)

	for i := 0; i < s.workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range pathChan {
				if err := s.processFile(path); err != nil {
					log.Error().Msgf("%s", err)
				}
			}
		}()
	}

	for _, path := range filePaths {
		pathChan <- path
	}
	close(pathChan)

	wg.Wait()
}

func (s *Scheduler) processFile(path string) error {
	fmt.Println("Processing file:", path)
	return nil
}
