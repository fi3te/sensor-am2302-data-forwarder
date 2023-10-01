package main

import (
	"strconv"
	"time"

	"github.com/fi3te/sensor-am2302-data-forwarder/pkg/config"
	"github.com/fi3te/sensor-am2302-data-forwarder/pkg/control"
	"github.com/fi3te/sensor-am2302-data-forwarder/pkg/destination"
	"github.com/fi3te/sensor-am2302-data-forwarder/pkg/domain"
	"github.com/fi3te/sensor-am2302-data-forwarder/pkg/io"
	"github.com/fi3te/sensor-am2302-data-forwarder/pkg/logging"
)

const (
	expectedLineSize     int = 48
	additionalCharacters int = 4
	charactersToRead     int = expectedLineSize + additionalCharacters
)

var cfg *config.AppConfig
var forwarder destination.Forwarder

func main() {
	ls := logging.New(logging.LevelInfo)

	ls.Info.Println("Reading configuration...")
	var err error
	cfg, err = config.ReadConfig()
	if err != nil {
		panic(err)
	}
	interval := cfg.Interval()
	ls.Debug.Println("Interval: " + interval.String())
	ls.Debug.Println("First retry after error: " + cfg.FirstRetryAfterError().String())
	ls.Debug.Println("Source directory: " + cfg.SourceDirectory)
	ls.Debug.Println("File determination by date: " + strconv.FormatBool(cfg.FileDeterminationByDate))
	ls.Debug.Println("Retention period: " + cfg.RetentionPeriod().String())

	forwarder, err = destination.NewForwarder(cfg)
	if err != nil {
		panic(err)
	}

	stopTickerChan := make(chan bool)

	ls.Info.Println("Starting forwarder...")
	go control.RunTicker(interval, stopTickerChan, task, ls)

	control.WaitForInterrupt()

	ls.Info.Println("Stopping application...")
	stopTickerChan <- true
}

func task(ticker *control.StatefulTicker, t time.Time, ls *logging.LogSetup) {
	ls.Debug.Println("<<")
	err := forwardData(ls)

	interval := ticker.CurrentInterval
	if err != nil {
		ls.Error.Println(err)
		var newInterval time.Duration
		if interval == ticker.InitialInterval {
			newInterval = cfg.FirstRetryAfterError()
		} else {
			newInterval = interval * 2
		}
		ls.Debug.Printf("Waiting for %s after error...\n", newInterval)
		ticker.Reset(newInterval)
	} else if interval != ticker.InitialInterval {
		ls.Debug.Println("Resetting interval to initial value...")
		ticker.ResetToInitialInterval()
	}
	ls.Debug.Println(">>")
}

func forwardData(ls *logging.LogSetup) error {
	file, err := determineFile()
	if err != nil {
		return err
	}

	line, err := io.ReadLastLine(file.FilePath, charactersToRead)
	if err != nil {
		return err
	}
	ls.Debug.Printf("Data: %s\n", line)

	dataPoint, err := domain.Parse(line)
	if err != nil {
		return err
	}

	date := file.FileNameWithoutExtension()
	ttl := int(time.Now().Add(cfg.RetentionPeriod()).Unix())
	ls.Debug.Println("Forwarding data...")
	err = forwarder.Forward(date, ttl, dataPoint)

	return err
}

func determineFile() (*io.File, error) {
	if cfg.FileDeterminationByDate {
		return io.DetermineFileByDate(cfg.SourceDirectory)
	} else {
		return io.DetermineFileByOrder(cfg.SourceDirectory)
	}
}
