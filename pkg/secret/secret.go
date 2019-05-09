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
package secret

import (
	"errors"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"

	anno "github.com/cmattoon/aws-ssm/pkg/annotations"
	"github.com/cmattoon/aws-ssm/pkg/provider"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

type Secret struct {
	Secret v1.Secret
	// Kubernetes Secret Name
	Name string
	// Kubernetes Namespace
	Namespace string
	// AWS Param Name
	ParamName string
	// AWS Param Type
	ParamType string
	// AWS Param Key (Default: "alias/aws/ssm")
	ParamKey string
	// AWS Param Value
	ParamValue string
	// The data to add to Kubernetes Secret Data
	Data map[string]string
}

func NewSecret(sec v1.Secret, p provider.Provider, secret_name string, secret_namespace string, param_name string, param_type string, param_key string) (*Secret, error) {

	s := &Secret{
		Secret:     sec,
		Name:       secret_name,
		Namespace:  secret_namespace,
		ParamName:  param_name,
		ParamType:  param_type,
		ParamKey:   param_key,
		ParamValue: "",
		Data:       map[string]string{},
	}

	log.Debugf("Getting value for '%s/%s'", s.Namespace, s.Name)

	decrypt := false
	if s.ParamKey != "" {
		decrypt = true
	}

	if s.ParamType == "String" || s.ParamType == "SecureString" {
		value, err := p.GetParameterValue(s.ParamName, decrypt)
		if err != nil {
			return nil, err
		}
		s.ParamValue = value
	} else if s.ParamType == "StringList" {
		value, err := p.GetParameterValue(s.ParamName, decrypt)
		if err != nil {
			return nil, err
		}
		s.ParamValue = value
		// StringList: Also set each key
		values := s.ParseStringList()
		for k, v := range values {
			s.Set(k, v)
		}
	} else if s.ParamType == "Directory" {
		// Directory: Set each sub-key
		all_params, err := p.GetParameterDataByPath(s.ParamName, decrypt)
		if err != nil {
			return nil, err
		}

		for k, v := range all_params {
			s.Set(safeKeyName(k), v)
		}
		s.ParamValue = "true" // Reads "Directory": "true"
		return s, nil
	}

	// Always set the "$ParamType" key:
	//   String: Value
	//   SecureString: Value
	//   StringList: Value
	//   Directory: <ssm-path>
	s.Set(s.ParamType, s.ParamValue)

	return s, nil
}

// FromKubernetesSecret returns an internal Secret struct, if the v1.Secret is properly annotated.
func FromKubernetesSecret(p provider.Provider, secret v1.Secret) (*Secret, error) {
	param_name := ""
	param_type := ""
	param_key := ""

	for k, v := range secret.ObjectMeta.Annotations {
		switch k {
		case anno.AWSParamName, anno.V1ParamName:
			param_name = v
		case anno.AWSParamType, anno.V1ParamType:
			param_type = v
		case anno.AWSParamKey, anno.V1ParamKey:
			param_key = v
		}
	}

	if param_name == "" || param_type == "" {
		return nil, errors.New("Irrelevant Secret")
	}

	if param_name != "" && param_type != "" {
		if param_type == "SecureString" && param_key == "" {
			log.Info("No KMS key defined. Using default key 'alias/aws/ssm'")
			param_key = "alias/aws/ssm"
		}
	}

	s, err := NewSecret(
		secret,
		p,
		secret.ObjectMeta.Name,
		secret.ObjectMeta.Namespace,
		param_name,
		param_type,
		param_key)

	if err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Secret) ParseStringList() (values map[string]string) {
	values = make(map[string]string)

	for _, pair := range strings.Split(strings.TrimSpace(s.ParamValue), ",") {
		pair = strings.TrimSpace(pair)
		key := pair
		val := ""

		if strings.Contains(pair, "=") {
			kv := strings.SplitN(pair, "=", 2)
			if len(kv) == 2 {
				if kv[0] != "" {
					key = kv[0]
					val = kv[1]
				}
			}
		}
		if key != "" {
			values[key] = val
		}
	}

	return
}

func (s *Secret) Set(key string, val string) (err error) {
	log.Debugf("Setting key=%s", key)
	if s.Secret.StringData == nil {
		s.Secret.StringData = make(map[string]string)
	}
	// StringData isn't populated initially, so check s.Data
	if _, ok := s.Data[key]; ok {
		// Refuse to overwite existing keys
		return errors.New(fmt.Sprintf("Key '%s' already exists for Secret %s/%s", key, s.Namespace, s.Name))
	}
	s.Secret.StringData[key] = val
	return
}

func (s *Secret) UpdateObject(cli kubernetes.Interface) (result *v1.Secret, err error) {
	log.Info("Updating Kubernetes Secret...")
	return cli.CoreV1().Secrets(s.Namespace).Update(&s.Secret)
}

func safeKeyName(key string) string {
	key = strings.TrimRight(key, "/")
	if strings.HasPrefix(key, "/") {
		key = strings.Replace(key, "/", "", 1)
	}
	return strings.Replace(key, "/", "_", -1)
}
