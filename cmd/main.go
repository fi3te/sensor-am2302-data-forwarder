package main

import (
	"log"
	"strconv"
	"time"

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
	log.Println("FileDeterminationByDate: " + strconv.FormatBool(cfg.FileDeterminationByDate))
	log.Println("Destination url: " + cfg.DestinationUrl)

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
			ticker.Reset(time.Duration(cfg.FirstRetryAfterErrorInSeconds * int(time.Second)))
		} else {
			ticker.Reset(interval * 2)
		}

	} else if interval != ticker.InitialInterval {
		ticker.ResetToInitialInterval()
	}
}

func forwardData() error {
	filePath, err := determineFilePath()
	if err != nil {
		return err
	}

	content, err := io.ReadLastLine(filePath, charactersToRead)
	if err != nil {
		return err
	}

	dataPoint, err := domain.Parse(content)
	if err != nil {
		return err
	}

	_, err = io.HttpPostJson(cfg.DestinationUrl, *dataPoint)

	return nil
}

func determineFilePath() (string, error) {
	if cfg.FileDeterminationByDate {
		return io.DetermineFilePathByDate(cfg.SourceDirectory)
	} else {
		return io.DetermineFilePathByOrder(cfg.SourceDirectory)
	}
}
