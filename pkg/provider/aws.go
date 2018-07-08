package provider

import (
	"github.com/cmattoon/aws-ssm/pkg/config"
)

type AWSProvider struct {
	
}

func NewAWSProvider(cfg *config.Config) (Provider, error) {
	return AWSProvider{}, nil
}

func (p AWSProvider) GetParameter(name string) (string) {
	return ""
}
