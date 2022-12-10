package domain

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

func Parse(line string) (*DataPoint, error) {
	if len(line) < 37 {
		return nil, errors.New("Line is shorter than expected!")
	}

	components := strings.Fields(line)
	if len(components) != 5 {
		return nil, errors.New("Invalid data format!")
	}

	timeFormat, _ := regexp.Compile("(0[0-9]|1[0-9]|2[0-3])[:][0-5][0-9][:][0-5][0-9]")
	time := components[0]
	timeValid := timeFormat.MatchString(time)
	if !timeValid {
		return nil, errors.New("Invalid time format: " + time)
	}

	temperature := strings.TrimSuffix(components[2], "°C,")
	temperatureValue, err := strconv.ParseFloat(temperature, 64)
	if err != nil {
		return nil, errors.New("Invalid temperature format: " + temperature)
	}

	humidity := strings.TrimSuffix(components[4], "%")
	humidityValue, err := strconv.ParseFloat(humidity, 64)
	if err != nil {
		return nil, errors.New("Invalid humidity format: " + humidity)
	}

	dataPoint := DataPoint{
		Time:        time,
		Temperature: temperatureValue,
		Humidity:    humidityValue,
	}

	return &dataPoint, nil
}
