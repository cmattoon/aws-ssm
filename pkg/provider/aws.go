/**
 * Copyright 2018 Curtis Mattoon
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package provider

import (
	"path"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/cmattoon/aws-ssm/pkg/config"
	log "github.com/sirupsen/logrus"
)

type AWSProvider struct {
	Session *session.Session
	Service *ssm.SSM
}

func NewAWSProvider(cfg *config.Config) (Provider, error) {
	log.Info("ENTER NewAWSProvider: creating new session")
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(cfg.AWSRegion),
	})

	if err != nil {
		log.Fatalf("%s", err)
	}

	return AWSProvider{
		Session: sess,
		Service: ssm.New(sess),
	}, nil
}

func (p AWSProvider) GetParameterValue(name string, decrypt bool, roleArn string) (string, error) {
	param, err := p.Service.GetParameter(&ssm.GetParameterInput{
		Name:           aws.String(name),
		WithDecryption: aws.Bool(decrypt),
	})

	if err != nil {
		log.Errorf("Failed to GetParameterValue: %s", err)
		return "", err
	}

	return *param.Parameter.Value, nil
}

func (p AWSProvider) GetParameterDataByPath(ppath string, decrypt bool, roleArn string) (map[string]string, error) {
	log.Info("ENTER GetParameterDataByPath with RoleArn: ", roleArn)
	if roleArn != "" {
		creds := stscreds.NewCredentials(p.Session, roleArn)
		p.Service = ssm.New(p.Session, &aws.Config{
			Credentials: creds,
			Region:      aws.String(*p.Session.Config.Region),
		})
	}
	// ppath is something like /path/to/env
	params, err := p.Service.GetParametersByPath(&ssm.GetParametersByPathInput{
		Path:           aws.String(ppath),
		Recursive:      aws.Bool(true),
		WithDecryption: aws.Bool(decrypt),
	})

	if err != nil {
		log.Errorf("Failed to GetParameterDataByPath: %s", err)
		return nil, err
	}

	results := make(map[string]string)
	// '/path/to/env/foo' -> 'foo': *pa.Value
	for _, pa := range params.Parameters {
		_, basename := path.Split(*pa.Name)
		results[basename] = *pa.Value
	}
	return results, nil
}
