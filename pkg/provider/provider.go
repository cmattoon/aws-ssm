package provider

import (
	//log "github.com/sirupsen/logrus"
	"github.com/cmattoon/aws-ssm/pkg/config"
)


type Provider interface {
	GetParameter(string)(string)
}


func NewProvider(cfg *config.Config) (Provider, error) {
	p, err := NewAWSProvider(cfg)
	// switch cfg.Provider {
	// case "aws":
	// 	log.Info("Using provider: aws")
	// }
	return p, err
}
