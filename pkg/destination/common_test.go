package destination

import (
	"testing"

	"github.com/fi3te/sensor-am2302-data-forwarder/pkg/domain"
	"github.com/stretchr/testify/assert"
)

func TestEncodeDataPointAsJson(t *testing.T) {
	dataPoint := domain.DataPoint{Time: "00:00:00", Temperature: 20.0, Humidity: 50.0}

	body, err := encodeDataPointAsJson("2023-09-26", 0, &dataPoint)

	assert.Nil(t, err)
	assert.Equal(t, "{\"date\":\"2023-09-26\",\"time\":\"00:00:00\",\"temperature\":20,\"humidity\":50,\"ttl\":0}\n", body.String())
}
