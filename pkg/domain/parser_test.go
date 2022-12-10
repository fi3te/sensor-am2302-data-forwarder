package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseValid(t *testing.T) {
	line := "00:00:00 temperature: 18.20°C, humidity: 61.70%"

	dataPoint, err := Parse(line)

	assert.Nil(t, err)
	assert.Equal(t, "00:00:00", dataPoint.Time)
	assert.Equal(t, 18.2, dataPoint.Temperature)
	assert.Equal(t, 61.7, dataPoint.Humidity)
}

func TestParseInvalidTime(t *testing.T) {
	line := "24:00:00 temperature: 18.20°C, humidity: 61.70%"

	dataPoint, err := Parse(line)

	assert.Nil(t, dataPoint)
	assert.NotNil(t, err)
}

func TestParseInvalidTemperature(t *testing.T) {
	line := "00:00:00 temperature: a18.20°C, humidity: 61.70%"

	dataPoint, err := Parse(line)

	assert.Nil(t, dataPoint)
	assert.NotNil(t, err)
}

func TestParseInvalidHumidity(t *testing.T) {
	line := "00:00:00 temperature: 18.20°C, humidity: 6+1.70%"

	dataPoint, err := Parse(line)

	assert.Nil(t, dataPoint)
	assert.NotNil(t, err)
}
