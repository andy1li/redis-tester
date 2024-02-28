package internal

import (
	"github.com/codecrafters-io/redis-tester/internal/command_test"
	"github.com/codecrafters-io/redis-tester/internal/instrumented_redis_client"
	"github.com/codecrafters-io/redis-tester/internal/resp_assertions"
	testerutils "github.com/codecrafters-io/tester-utils"
	"github.com/codecrafters-io/tester-utils/random"
)

// Tests 'GET, SET'
func testGetSet(stageHarness *testerutils.StageHarness) error {
	b := NewRedisBinary(stageHarness)
	if err := b.Run(); err != nil {
		return err
	}

	logger := stageHarness.Logger

	client, err := instrumented_redis_client.NewInstrumentedRedisClient(stageHarness, "localhost:6379", "")
	if err != nil {
		logFriendlyError(logger, err)
		return err
	}

	randomWords := random.RandomWords(2)

	randomKey := randomWords[0]
	randomValue := randomWords[1]

	logger.Debugf("Setting key %s to %s", randomKey, randomValue)
	setCommandTestCase := command_test.CommandTestCase{
		Command:   "set",
		Args:      []string{randomKey, randomValue},
		Assertion: resp_assertions.NewStringValueAssertion("OK"),
	}

	if err := setCommandTestCase.Run(client, logger); err != nil {
		logFriendlyError(logger, err)
		return err
	}

	logger.Debugf("Getting key %s", randomKey)

	getCommandTestCase := command_test.CommandTestCase{
		Command:   "get",
		Args:      []string{randomKey},
		Assertion: resp_assertions.NewStringValueAssertion(randomValue),
	}

	if err := getCommandTestCase.Run(client, logger); err != nil {
		logFriendlyError(logger, err)
		return err
	}

	client.Close()
	return nil
}
