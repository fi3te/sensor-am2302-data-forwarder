package destination

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/config"
	appConfig "github.com/fi3te/sensor-am2302-data-forwarder/pkg/config"
	"github.com/fi3te/sensor-am2302-data-forwarder/pkg/domain"
)

type AwsForwarder struct {
	url              string
	method           string
	expectedStatus   int
	cfg              *aws.Config
	credentialsCache *aws.CredentialsCache
}

func NewAwsForwarder(appConfig *appConfig.AwsConfig) (*AwsForwarder, error) {
	awsConfig, err := config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile(appConfig.Profile))
	if err != nil {
		return nil, err
	}
	credentialsCache := aws.NewCredentialsCache(awsConfig.Credentials)
	return &AwsForwarder{appConfig.Url, appConfig.Method, appConfig.ExpectedStatus, &awsConfig, credentialsCache}, err
}

func (f *AwsForwarder) Forward(date string, ttl int64, dataPoint *domain.DataPoint) error {
	req, err := f.buildRequest(date, ttl, dataPoint)
	if err != nil {
		return err
	}
	return sendRequest(req, f.expectedStatus)
}

func (f *AwsForwarder) buildRequest(date string, ttl int64, dataPoint *domain.DataPoint) (*http.Request, error) {
	credentials, err := f.credentialsCache.Retrieve(context.TODO())
	if err != nil {
		return nil, err
	}

	req, body, err := buildRequest(f.method, f.url, date, ttl, dataPoint)
	if err != nil {
		return nil, err
	}

	hash := sha256.New()
	hash.Write(body)
	encodedHash := hex.EncodeToString(hash.Sum(nil))

	signer := v4.NewSigner()
	err = signer.SignHTTP(context.TODO(), credentials, req, encodedHash, "execute-api", f.cfg.Region, time.Now())
	if err != nil {
		return nil, err
	}

	return req, nil
}
