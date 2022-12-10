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
	DestinationUrl                string `yaml:"destination-url"`
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
	if cfg.DestinationUrl == "" {
		return errors.New("Attribute 'Destination url' must be specified.")
	}
	return nil
}

func (cfg *Config) Interval() time.Duration {
	return secondsToDuration(cfg.IntervalInSeconds)
}

func (cfg *Config) FirstRetryAfterError() time.Duration {
	return secondsToDuration(cfg.FirstRetryAfterErrorInSeconds)
}

func secondsToDuration(seconds int) time.Duration {
	return time.Duration(seconds * int(time.Second))
}
