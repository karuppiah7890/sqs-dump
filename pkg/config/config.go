package config

import (
	"fmt"
	"net/url"
	"os"
)

// All configuration is through environment variables

const AWS_REGION_ENV_VAR = "AWS_REGION"
const AWS_ACCESS_KEY_ID_ENV_VAR = "AWS_ACCESS_KEY_ID"
const AWS_SECRET_ACCESS_KEY_ENV_VAR = "AWS_SECRET_ACCESS_KEY"
const SQS_QUEUE_URL_ENV_VAR = "SQS_QUEUE_URL"

type Config struct {
	awsRegion          string
	awsAccessKeyId     string
	awsSecretAccessKey string
	sqsQueueUrl        string
}

func NewConfigFromEnvVars() (*Config, error) {
	awsRegion, err := getAwsRegion()
	if err != nil {
		return nil, fmt.Errorf("error occurred while getting AWS Region: %v", err)
	}

	awsAccessKeyId, err := getAwsAccessKeyId()
	if err != nil {
		return nil, fmt.Errorf("error occurred while getting AWS Access Key ID: %v", err)
	}

	awsSecretAccessKey, err := getAwsSecretAccessKey()
	if err != nil {
		return nil, fmt.Errorf("error occurred while getting AWS Secret Access Key: %v", err)
	}

	sqsQueueUrl, err := getSqsQueueUrl()
	if err != nil {
		return nil, fmt.Errorf("error occurred while getting SQS Queue URL: %v", err)
	}

	return &Config{
		awsRegion:          awsRegion,
		awsAccessKeyId:     awsAccessKeyId,
		awsSecretAccessKey: awsSecretAccessKey,
		sqsQueueUrl:        sqsQueueUrl,
	}, nil
}

// Get SQS Queue URL
func getSqsQueueUrl() (string, error) {
	sqsQueueUrl, ok := os.LookupEnv(SQS_QUEUE_URL_ENV_VAR)
	if !ok {
		return "", fmt.Errorf("%s environment variable value is a required value. Please provide it", SQS_QUEUE_URL_ENV_VAR)
	}

	_, err := url.Parse(sqsQueueUrl)
	if err != nil {
		return "", fmt.Errorf("error while parsing %s environment variable value (%s): %v", SQS_QUEUE_URL_ENV_VAR, sqsQueueUrl, err)
	}

	return sqsQueueUrl, nil
}

// Get AWS Region
func getAwsRegion() (string, error) {
	awsRegion, ok := os.LookupEnv(AWS_REGION_ENV_VAR)
	if !ok {
		return "", fmt.Errorf("%s environment variable value is a required value. Please define it", AWS_ACCESS_KEY_ID_ENV_VAR)
	}

	return awsRegion, nil
}

// Get AWS Access Key ID
func getAwsAccessKeyId() (string, error) {
	awsAccessKeyId, ok := os.LookupEnv(AWS_ACCESS_KEY_ID_ENV_VAR)
	if !ok {
		return "", fmt.Errorf("%s environment variable value is a required value. Please define it", AWS_ACCESS_KEY_ID_ENV_VAR)
	}

	return awsAccessKeyId, nil
}

// Get AWS Secret Access Key
func getAwsSecretAccessKey() (string, error) {
	awsSecretAccessKey, ok := os.LookupEnv(AWS_SECRET_ACCESS_KEY_ENV_VAR)
	if !ok {
		return "", fmt.Errorf("%s environment variable value is a required value. Please define it", AWS_SECRET_ACCESS_KEY_ENV_VAR)
	}

	return awsSecretAccessKey, nil
}

func (c *Config) GetSqsQueueUrl() string {
	return c.sqsQueueUrl
}
