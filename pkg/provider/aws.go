package provider

import (
	log "github.com/sirupsen/logrus"
	
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/cmattoon/aws-ssm/pkg/config"
)

type AWSProvider struct {
	Session *session.Session
	Service *ssm.SSM
}

func NewAWSProvider(cfg *config.Config) (Provider, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(cfg.AWSRegion),
		Credentials: credentials.NewEnvCredentials(),
	})

	if err != nil {
		log.Fatalf("%s", err)
	}
	
	return AWSProvider{
		Session: sess,
		Service: ssm.New(sess),
	}, nil
}

func (p AWSProvider) GetParameterValue(name string, decrypt bool) (string, error) {
	log.Debugf("GetParameterValue(%v, %v)", name, decrypt)
	param, err := p.Service.GetParameter(&ssm.GetParameterInput{
		Name: aws.String(name),
		WithDecryption: aws.Bool(decrypt),
	})
	
	if err != nil {
		log.Debug("Failed to get value. Returning ''")
		return "", err
	}
	
	return *param.Parameter.Value, nil
}
