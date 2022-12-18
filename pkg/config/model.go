package config

import (
	"errors"
	"time"
)

type Config struct {
	IntervalInSeconds             int    `yaml:"interval-in-seconds"`
	FirstRetryAfterErrorInSeconds int    `yaml:"first-retry-after-error-in-seconds"`
	SourceDirectory               string `yaml:"source-directory"`
	FileDeterminationByDate       bool   `yaml:"file-determination-by-date"`
	AwsProfile                    string `yaml:"aws-profile"`
	AwsApiGatewayDestinationUrl   string `yaml:"aws-api-gateway-destination-url"`
	AwsRetentionPeriodInHours     int    `yaml:"aws-retention-period-in-hours"`
}

func (cfg *Config) validate() error {
	if cfg.IntervalInSeconds < 1 {
		return errors.New("Attribute 'Interval in seconds' must be greater than zero.")
	}
	if cfg.FirstRetryAfterErrorInSeconds < 1 {
		return errors.New("Attribute 'First retry after error in seconds' must be greater than zero.")
	}
	if cfg.SourceDirectory == "" {
		return errors.New("Attribute 'Source directory' must be specified.")
	}
	if cfg.AwsProfile == "" {
		return errors.New("Attribute 'AWS profile' must be specified.")
	}
	if cfg.AwsApiGatewayDestinationUrl == "" {
		return errors.New("Attribute 'AWS API Gateway destination url' must be specified.")
	}
	if cfg.AwsRetentionPeriodInHours < 1 {
		return errors.New("Attribute 'AWS retention period in hours' must be greater than zero.")
	}
	return nil
}

func (cfg *Config) Interval() time.Duration {
	return intToDuration(cfg.IntervalInSeconds, time.Second)
}

func (cfg *Config) FirstRetryAfterError() time.Duration {
	return intToDuration(cfg.FirstRetryAfterErrorInSeconds, time.Second)
}

func (cfg *Config) RetentionPeriod() time.Duration {
	return intToDuration(cfg.AwsRetentionPeriodInHours, time.Hour)
}

func intToDuration(value int, unit time.Duration) time.Duration {
	return time.Duration(value * int(unit))
}
