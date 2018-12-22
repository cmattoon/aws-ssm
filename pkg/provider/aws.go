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
	log "github.com/sirupsen/logrus"
	"strings"

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
		Region:      aws.String(cfg.AWSRegion),
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
		Name:           aws.String(name),
		WithDecryption: aws.Bool(decrypt),
	})

	if err != nil {
		log.Debug("Failed to get value. Returning ''")
		return "", err
	}

	return *param.Parameter.Value, nil
}

func (p AWSProvider) GetParameterDataByPath(path string, decrypt bool) (map[string]string, error) {
	// path is something like /path/to/env
	params, err := p.Service.GetParametersByPath(&ssm.GetParametersByPathInput{
		Path:           aws.String(path),
		Recursive:      aws.Bool(true),
		WithDecryption: aws.Bool(decrypt),
	})

	if err != nil {
		log.Fatal("Failed to get params by path '%s': %s", path, err)
	}

	results := make(map[string]string)
	// '/path/to/env/foo' -> 'foo': *pa.Value
	for _, pa := range params.Parameters {
		basename := strings.Replace(*pa.Name, path, "", -1)
		results[basename] = *pa.Value
	}
	return results, nil
}
