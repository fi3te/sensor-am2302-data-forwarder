package config

import (
	"fmt"
	"net/http"
	"time"
)

type HttpConfig struct {
	Url            string `yaml:"url"`
	Method         string `yaml:"method"`
	ExpectedStatus int    `yaml:"expected-status"`
}

func (cfg *HttpConfig) validate() error {
	if cfg.Url == "" {
		return errRequired("http url")
	}
	if cfg.Method == "" {
		return errRequired("http method")
	}
	if cfg.ExpectedStatus < 200 || cfg.ExpectedStatus > 300 {
		return errBetween("expected status", 200, 299)
	}
	return nil
}

type PlainHttpConfig struct {
	HttpConfig        `yaml:",inline"`
	BasicAuthUsername string `yaml:"basic-auth-username"`
	BasicAuthPassword string `yaml:"basic-auth-password"`
}

func (cfg *PlainHttpConfig) validate() error {
	err := cfg.HttpConfig.validate()
	if err != nil {
		return err
	}
	_, err = cfg.BuildAuthentication()
	return err
}

type HttpAuthentication interface {
	Apply(req *http.Request)
}

type noAuthentication struct{}

func (na *noAuthentication) Apply(req *http.Request) {}

type basicAuthentication struct {
	username string
	password string
}

func (ba *basicAuthentication) Apply(req *http.Request) {
	req.SetBasicAuth(ba.username, ba.password)
}

func (cfg *PlainHttpConfig) BuildAuthentication() (HttpAuthentication, error) {
	if cfg.BasicAuthUsername != "" || cfg.BasicAuthPassword != "" {
		if cfg.BasicAuthUsername == "" {
			return nil, errRequired("basic authentication username")
		}
		if cfg.BasicAuthPassword == "" {
			return nil, errRequired("basic authentication password")
		}
		return &basicAuthentication{username: cfg.BasicAuthUsername, password: cfg.BasicAuthPassword}, nil
	}
	return &noAuthentication{}, nil
}

type AwsConfig struct {
	HttpConfig `yaml:",inline"`
	Profile    string `yaml:"profile"`
}

func (cfg *AwsConfig) validate() error {
	err := cfg.HttpConfig.validate()
	if err != nil {
		return err
	}
	if cfg.Profile == "" {
		return errRequired("AWS profile")
	}
	return nil
}

type NtfyConfig struct {
	Url             string   `yaml:"url"`
	Topic           string   `yaml:"topic"`
	TitleTemplate   string   `yaml:"title-template"`
	MessageTemplate string   `yaml:"message-template"`
	Tags            []string `yaml:"tags"`
}

func (cfg *NtfyConfig) validate() error {
	if cfg.Url == "" {
		return errRequired("ntfy url")
	}
	if cfg.Topic == "" {
		return errRequired("ntfy topic")
	}
	if cfg.TitleTemplate == "" {
		return errRequired("ntfy title template")
	}
	if cfg.MessageTemplate == "" {
		return errRequired("ntfy message template")
	}
	return nil
}

type AppConfig struct {
	IntervalInSeconds             int64            `yaml:"interval-in-seconds"`
	FirstRetryAfterErrorInSeconds int64            `yaml:"first-retry-after-error-in-seconds"`
	SourceDirectory               string           `yaml:"source-directory"`
	FileDeterminationByDate       bool             `yaml:"file-determination-by-date"`
	RetentionPeriodInHours        int64            `yaml:"retention-period-in-hours"`
	Aws                           *AwsConfig       `yaml:"aws"`
	Http                          *PlainHttpConfig `yaml:"http"`
	Ntfy                          *NtfyConfig      `yaml:"ntfy"`
}

func (cfg *AppConfig) validate() error {
	if cfg.IntervalInSeconds < 1 {
		return errGtZero("interval in seconds")
	}
	if cfg.FirstRetryAfterErrorInSeconds < 1 {
		return errGtZero("first retry after error in seconds")
	}
	if cfg.SourceDirectory == "" {
		return errRequired("source directory")
	}
	if cfg.RetentionPeriodInHours < 1 {
		return errGtZero("retention period in hours")
	}

	destinations := count(cfg.Aws != nil, cfg.Http != nil, cfg.Ntfy != nil)
	if destinations != 1 {
		return fmt.Errorf("only exactly one destination at a time is supported (specified: %d)", destinations)
	}
	if cfg.Aws != nil {
		return cfg.Aws.validate()
	}
	if cfg.Http != nil {
		return cfg.Http.validate()
	}
	return cfg.Ntfy.validate()
}

func (cfg *AppConfig) Interval() time.Duration {
	return intToDuration(cfg.IntervalInSeconds, time.Second)
}

func (cfg *AppConfig) FirstRetryAfterError() time.Duration {
	return intToDuration(cfg.FirstRetryAfterErrorInSeconds, time.Second)
}

func (cfg *AppConfig) RetentionPeriod() time.Duration {
	return intToDuration(cfg.RetentionPeriodInHours, time.Hour)
}

func intToDuration(value int64, unit time.Duration) time.Duration {
	return time.Duration(value * int64(unit))
}

func errRequired(description string) error {
	return fmt.Errorf("attribute '%s' must be specified", description)
}

func errGtZero(description string) error {
	return fmt.Errorf("attribute '%s' must be greater than zero", description)
}

func errBetween(description string, min int, max int) error {
	return fmt.Errorf("attribute '%s' must be between %d and %d", description, min, max)
}

func count(values ...bool) int {
	var count int
	for _, value := range values {
		if value {
			count++
		}
	}
	return count
}
