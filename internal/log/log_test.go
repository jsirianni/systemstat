package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConsts(t *testing.T) {
	// if these values are changed, we need to update docker-compose,
	// documentation, and any deployment configuration
	assert.Equal(t, "SYSTEMSTAT_LOG_LEVEL", envLogLevel)
	assert.Equal(t, "ERROR", errorLVL)
	assert.Equal(t, "INFO", infoLVL)
	assert.Equal(t, "DEBUG", debugLVL)
	assert.Equal(t, "TRACE", traceLVL)
	assert.Equal(t, 1, errorLevel)
	assert.Equal(t, 2, infoLevel)
	assert.Equal(t, 3, debugLevel)
	assert.Equal(t, 4, traceLevel)
}

func TestSetLogLevel(t *testing.T) {
	levelMap := make(map[string]int)
	levelMap["ERROR"] = 1
	levelMap["INFO"] = 2
	levelMap["DEBUG"] = 3
	levelMap["TRACE"] = 4
	for k, v := range levelMap {
		SetLogLevel(k)
		assert.Equal(t, v, level)
	}
}
