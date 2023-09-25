package destination

import (
	"bytes"
	"encoding/json"
	"errors"

	"github.com/fi3te/sensor-am2302-data-forwarder/pkg/config"
	"github.com/fi3te/sensor-am2302-data-forwarder/pkg/domain"
)

type Forwarder interface {
	Forward(date string, ttl int, dataPoint *domain.DataPoint) error
}

func NewForwarder(appConfig *config.AppConfig) (Forwarder, error) {
	if appConfig.Aws != nil {
		return NewAwsForwarder(appConfig.Aws)
	}
	if appConfig.Http != nil {
		return NewHttpForwarder(appConfig.Http)
	}
	if appConfig.Ntfy != nil {
		return NewNtfyForwarder(appConfig.Ntfy)
	}
	return nil, errors.New("no forwarder configured")
}

type dataPointDto struct {
	Date        string  `json:"date"`
	Time        string  `json:"time"`
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
	Ttl         int     `json:"ttl"`
}

func buildDataPointDto(date string, ttl int, dataPoint *domain.DataPoint) dataPointDto {
	return dataPointDto{
		Date:        date,
		Time:        dataPoint.Time,
		Temperature: dataPoint.Temperature,
		Humidity:    dataPoint.Humidity,
		Ttl:         ttl,
	}
}

func encodeDataPointAsJson(date string, ttl int, dataPoint *domain.DataPoint) (*bytes.Buffer, error) {
	dto := buildDataPointDto(date, ttl, dataPoint)
	return encodeAsJson(dto)
}

func encodeAsJson(v any) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(v)
	if err != nil {
		return nil, err
	}
	return &buf, nil
}
