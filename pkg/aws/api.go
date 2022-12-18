package aws

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/fi3te/sensor-am2302-data-forwarder/pkg/domain"
)

func InitConfig(awsProfile string, awsApiGatewayDestinationUrl string) (*AwsConfig, error) {
	awsConfig, err := config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile(awsProfile))
	if err != nil {
		return nil, err
	}

	credentialsCache := aws.NewCredentialsCache(awsConfig.Credentials)

	return &AwsConfig{
		Config:                      &awsConfig,
		CredentialsCache:            credentialsCache,
		AwsApiGatewayDestinationUrl: awsApiGatewayDestinationUrl,
	}, nil
}

func SendToApiGateway(awsCfg *AwsConfig, date string, ttl int, dataPoint *domain.DataPoint) error {
	req, err := createRequest(awsCfg, date, ttl, dataPoint)
	if err != nil {
		return err
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		return errors.New(fmt.Sprintf("AWS status code: %d", response.StatusCode))
	}

	return nil
}

func createRequest(awsCfg *AwsConfig, date string, ttl int, dataPoint *domain.DataPoint) (*http.Request, error) {
	credentials, err := awsCfg.CredentialsCache.Retrieve(context.TODO())
	if err != nil {
		return nil, err
	}

	body := buildAwsDataPoint(date, ttl, dataPoint)
	buf, encodedHash, err := encode(&body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, awsCfg.AwsApiGatewayDestinationUrl, buf)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	signer := v4.NewSigner()
	err = signer.SignHTTP(context.TODO(), credentials, req, encodedHash, "execute-api", awsCfg.Config.Region, time.Now())
	if err != nil {
		return nil, err
	}

	return req, nil
}

func buildAwsDataPoint(date string, ttl int, dataPoint *domain.DataPoint) AwsDataPoint {
	return AwsDataPoint{
		Date:        date,
		Time:        dataPoint.Time,
		Temperature: dataPoint.Temperature,
		Humidity:    dataPoint.Humidity,
		Ttl:         ttl,
	}
}

func encode(body *AwsDataPoint) (*bytes.Buffer, string, error) {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(body)
	if err != nil {
		return nil, "", err
	}

	hash := sha256.New()
	hash.Write(buf.Bytes())
	encodedHash := hex.EncodeToString(hash.Sum(nil))

	return &buf, encodedHash, nil
}
