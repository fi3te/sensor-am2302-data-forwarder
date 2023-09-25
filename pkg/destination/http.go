package destination

import (
	"fmt"
	"net/http"

	appConfig "github.com/fi3te/sensor-am2302-data-forwarder/pkg/config"
	"github.com/fi3te/sensor-am2302-data-forwarder/pkg/domain"
)

type HttpForwarder struct {
	url            string
	method         string
	expectedStatus int
	authentication appConfig.HttpAuthentication
}

func NewHttpForwarder(appConfig *appConfig.PlainHttpConfig) (*HttpForwarder, error) {
	authentication, err := appConfig.BuildAuthentication()
	if err != nil {
		return nil, err
	}
	return &HttpForwarder{appConfig.Url, appConfig.Method, appConfig.ExpectedStatus, authentication}, nil
}

func (f *HttpForwarder) Forward(date string, ttl int, dataPoint *domain.DataPoint) error {
	req, _, err := buildRequest(f.method, f.url, date, ttl, dataPoint)
	if err != nil {
		return err
	}
	f.authentication.Apply(req)
	return sendRequest(req, f.expectedStatus)
}

// general functionality

func buildRequest(method string, url string, date string, ttl int, dataPoint *domain.DataPoint) (*http.Request, []byte, error) {
	buf, err := encodeDataPointAsJson(date, ttl, dataPoint)
	if err != nil {
		return nil, nil, err
	}
	body := buf.Bytes()
	req, err := http.NewRequest(method, url, buf)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	return req, body, nil
}

func sendRequest(req *http.Request, expectedStatus int) error {
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode != expectedStatus {
		return fmt.Errorf("unexpected http status: %d", response.StatusCode)
	}

	return nil
}
