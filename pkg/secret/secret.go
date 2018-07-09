package secret

import (
	"errors"
	
	log "github.com/sirupsen/logrus"

	"k8s.io/api/core/v1"

	"github.com/cmattoon/aws-ssm/pkg/provider"
	"github.com/cmattoon/aws-ssm/pkg/annotations"
)

type Secret struct {
	Name string
	Namespace string
	ParamValue string
	Data map[string]string
}


func NewSecret(p provider.Provider, name string, ns string, decrypt bool) (*Secret) {
	s := &Secret{
		Name: name,
		Namespace: ns,
		Data: map[string]string{},
	}
	value, err := p.GetParameterValue(name, decrypt)
	if err != nil {
		log.Infof("Couldn't get value for %s: %s", name, err)
	} else {
		s.ParamValue = value
	}
	return s
}

func FromKubernetesSecret(p provider.Provider, secret v1.Secret) (*Secret, error) {
	param_name := ""
	param_type := ""
	param_key := ""
	
	for k, v := range secret.ObjectMeta.Annotations {
		switch k {
		case annotations.AwsParamName:
			param_name = v
		case annotations.AwsParamType:
			param_type = v
		case annotations.AwsParamKey:
			param_key = v
		}
	}
	
	if param_name == "" || param_type == "" {
		return nil, errors.New("Irrelevant Secret")
	}
	
	decrypt := false
	if param_name != "" && param_type != "" {
		if param_type == "SecureString" && param_key == "" {
			log.Info("No KMS key defined. Using default key 'alias/aws/ssm'")
			param_key = "alias/aws/ssm"
		}
	}
	
	// If a KMS key is specified (aws-param-key), decrypt it
	if param_key != "" {
		decrypt = true
	}
		
	s := NewSecret(p, secret.ObjectMeta.Name, secret.ObjectMeta.Namespace, decrypt)

	// Set 'secretValue' accoring to the Parameter value
	for k, v := range secret.StringData {
		s.Data[k] = v

		if k == "secretValue" {
			log.Infof("Overwriting key %s for secret %s", k, secret.ObjectMeta.Name)
		}
	}
	
	s.Data["secretValue"] = s.ParamValue
	return s, nil
}
