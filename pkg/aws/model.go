package aws

import (
	"github.com/aws/aws-sdk-go-v2/aws"
)

type AwsDataPoint struct {
	Date        string  `json:"date"`
	Time        string  `json:"time"`
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
	Ttl         int     `json:"ttl"`
}

type AwsConfig struct {
	Config                      *aws.Config
	CredentialsCache            *aws.CredentialsCache
	AwsApiGatewayDestinationUrl string
}
