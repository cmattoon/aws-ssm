package provider

import (
	//log "github.com/sirupsen/logrus"
	"github.com/cmattoon/aws-ssm/pkg/config"
)


type Provider interface {
	GetParameterValue(string, bool) (string, error)
}


func NewProvider(cfg *config.Config) (Provider, error) {
	p, err := NewAWSProvider(cfg)
	return p, err
}
