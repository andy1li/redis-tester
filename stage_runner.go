package main

import "time"
import "fmt"

// StageRunnerResult is returned from StageRunner.Run()
type StageRunnerResult struct {
	failedAtStage Stage
	error         error
	logger        *customLogger
}

// IsSuccess says whether a StageRunnerResult was successful
// or not
func (res StageRunnerResult) IsSuccess() bool {
	return res.error == nil
}

// StageRunner is used to run multiple stages
type StageRunner struct {
	stages  []Stage
	isDebug bool
}

func newStageRunner(isDebug bool) StageRunner {
	return StageRunner{
		isDebug: isDebug,
		stages: []Stage{
			Stage{
				name:    "Stage 1: PING <-> PONG",
				logger:  getLogger(isDebug, "[stage-1] "),
				runFunc: runStage1,
			},
			Stage{
				name:    "Stage 2: Multiple Clients",
				logger:  getLogger(isDebug, "[stage-2] "),
				runFunc: runStage2,
			},
			Stage{
				name:    "Stage 3: ECHO... O... O...",
				logger:  getLogger(isDebug, "[stage-3] "),
				runFunc: runStage3,
			},
		},
	}
}

// Run tests in a specific StageRunner
func (r StageRunner) Run() StageRunnerResult {
	for _, stage := range r.stages {
		logger := stage.logger
		logger.Infof("Running test: %s", stage.name)

		stageResultChannel := make(chan error, 1)
		go func() {
			err := stage.runFunc(logger)
			stageResultChannel <- err
		}()

		var err error
		select {
		case stageErr := <-stageResultChannel:
			err = stageErr
		case <-time.After(1 * time.Second):
			err = fmt.Errorf("timed out, test exceeded 1 seconds")
		}

		if err != nil {
			reportTestError(err, r.isDebug, logger)
			return StageRunnerResult{
				failedAtStage: stage,
				error:         err,
			}
		}

		logger.Successf("Test passed.")
	}

	return StageRunnerResult{
		error: nil,
	}
}

func reportTestError(err error, isDebug bool, logger *customLogger) {
	logger.Errorf("%s", err)
	if isDebug {
		logger.Errorf("Test failed")
	} else {
		logger.Errorf("Test failed " +
			"(try using the --debug flag to see more output)")
	}
}

// Stage is blah
type Stage struct {
	name        string
	description string
	runFunc     func(logger *customLogger) error
	logger      *customLogger
}