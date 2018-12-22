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
	//log "github.com/sirupsen/logrus"
	"github.com/cmattoon/aws-ssm/pkg/config"
)

type Provider interface {
	GetParameterValue(string, bool) (string, error)
	GetParameterDataByPath(string, bool) (map[string]string, error)
}

func NewProvider(cfg *config.Config) (Provider, error) {
	p, err := NewAWSProvider(cfg)
	return p, err
}

type MockProvider struct {
	Value             string
	DecryptedValue    string
	DirectoryContents map[string]string
}

func (mp MockProvider) GetParameterValue(s string, b bool) (string, error) {
	if b {
		// Decrypt flag
		return mp.DecryptedValue, nil
	}
	return mp.Value, nil
}

func (mp MockProvider) GetParameterDataByPath(s string, b bool) (map[string]string, error) {
	return mp.DirectoryContents, nil
}
