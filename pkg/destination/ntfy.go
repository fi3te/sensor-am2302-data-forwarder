package destination

import (
	"fmt"
	"net/http"

	appConfig "github.com/fi3te/sensor-am2302-data-forwarder/pkg/config"
	"github.com/fi3te/sensor-am2302-data-forwarder/pkg/domain"
)

type ntfyMessage struct {
	Topic   string   `json:"topic"`
	Title   string   `json:"title"`
	Message string   `json:"message"`
	Tags    []string `json:"tags"`
}

type NtfyForwarder struct {
	url             string
	topic           string
	titleTemplate   string
	messageTemplate string
	tags            []string
}

func NewNtfyForwarder(appConfig *appConfig.NtfyConfig) (*NtfyForwarder, error) {
	return &NtfyForwarder{
		url:             appConfig.Url,
		topic:           appConfig.Topic,
		titleTemplate:   appConfig.TitleTemplate,
		messageTemplate: appConfig.MessageTemplate,
		tags:            appConfig.Tags,
	}, nil
}

func (f *NtfyForwarder) Forward(date string, ttl int, dataPoint *domain.DataPoint) error {
	buf, err := encodeAsJson(ntfyMessage{
		Topic:   f.topic,
		Title:   fmt.Sprintf(f.titleTemplate, dataPoint.Time[0:5]),
		Message: fmt.Sprintf(f.messageTemplate, dataPoint.Temperature, dataPoint.Humidity),
		Tags:    f.tags,
	})
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, f.url, buf)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	return sendRequest(req, http.StatusOK)
}
