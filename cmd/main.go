package main

import (
	"log"
	"strconv"
	"time"

	"github.com/fi3te/sensor-am2302-data-forwarder/pkg/aws"
	"github.com/fi3te/sensor-am2302-data-forwarder/pkg/config"
	"github.com/fi3te/sensor-am2302-data-forwarder/pkg/control"
	"github.com/fi3te/sensor-am2302-data-forwarder/pkg/domain"
	"github.com/fi3te/sensor-am2302-data-forwarder/pkg/io"
)

const (
	expectedLineSize     int = 48
	additionalCharacters int = 4
	charactersToRead     int = expectedLineSize + additionalCharacters
)

var cfg *config.Config
var awsCfg *aws.AwsConfig

func main() {
	log.Println("Reading configuration...")
	var err error
	cfg, err = config.ReadConfig()
	if err != nil {
		panic(err)
	}
	interval := cfg.Interval()
	log.Println("Interval: " + interval.String())
	log.Println("First retry after error: " + cfg.FirstRetryAfterError().String())
	log.Println("Source directory: " + cfg.SourceDirectory)
	log.Println("File determination by date: " + strconv.FormatBool(cfg.FileDeterminationByDate))
	log.Println("AWS profile: " + cfg.AwsProfile)
	log.Println("AWS API Gateway destination url: " + cfg.AwsApiGatewayDestinationUrl)
	log.Println("AWS retention period: " + cfg.RetentionPeriod().String())

	awsCfg, err = aws.InitConfig(cfg.AwsProfile, cfg.AwsApiGatewayDestinationUrl)
	if err != nil {
		panic(err)
	}

	stopTickerChan := make(chan bool)

	go control.RunTicker(interval, stopTickerChan, task)

	control.WaitForInterrupt()

	log.Println("Stopping application...")
	stopTickerChan <- true
}

func task(ticker *control.StatefulTicker, t time.Time) {
	err := forwardData()

	interval := ticker.CurrentInterval
	if err != nil {
		log.Println(err)
		if interval == ticker.InitialInterval {
			ticker.Reset(cfg.FirstRetryAfterError())
		} else {
			ticker.Reset(interval * 2)
		}

	} else if interval != ticker.InitialInterval {
		ticker.ResetToInitialInterval()
	}
}

func forwardData() error {
	file, err := determineFile()
	if err != nil {
		return err
	}

	line, err := io.ReadLastLine(file.FilePath, charactersToRead)
	if err != nil {
		return err
	}

	dataPoint, err := domain.Parse(line)
	if err != nil {
		return err
	}

	date := file.FileNameWithoutExtension()
	ttl := int(time.Now().Add(cfg.RetentionPeriod()).Unix())
	err = aws.SendToApiGateway(awsCfg, date, ttl, dataPoint)

	return err
}

func determineFile() (*io.File, error) {
	if cfg.FileDeterminationByDate {
		return io.DetermineFileByDate(cfg.SourceDirectory)
	} else {
		return io.DetermineFileByOrder(cfg.SourceDirectory)
	}
}
