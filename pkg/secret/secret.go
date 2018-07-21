package secret

import (
	"errors"
	"fmt"
	
	log "github.com/sirupsen/logrus"
	
	"k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"github.com/cmattoon/aws-ssm/pkg/provider"
	anno "github.com/cmattoon/aws-ssm/pkg/annotations"
	
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

func NewSecret(
	sec v1.Secret,
	p provider.Provider,
	secret_name string,
	secret_namespace string,
	param_name string,
	param_type string,
	param_key string,
) (*Secret) {
	
	s := &Secret{
		Secret: sec,
		Name: secret_name,
		Namespace: secret_namespace,
		ParamName: param_name,
		ParamType: param_type,
		ParamKey: param_key,
		ParamValue: "",
		Data: map[string]string{},
	}
	
	log.Infof("Getting value for '%s/%s'", s.Namespace, s.Name)

	decrypt := false
	if s.ParamKey != "" {
		decrypt = true
	}
	
	value, err := p.GetParameterValue(s.ParamName, decrypt)
	
	if err != nil {
		log.Infof("Couldn't get value for %s/%s: %s",
			s.Namespace, s.Name, err)
	} else {
		log.Info("Setting value to %v", value)
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
		case anno.AWSParamName:
			param_name = v
		case anno.AWSParamType:
			param_type = v
		case anno.AWSParamKey:
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
	
	s := NewSecret(
		secret,
		p,
		secret.ObjectMeta.Name,
		secret.ObjectMeta.Namespace,
		param_name,
		param_type,
		param_key)

	return s, nil
}

func (s *Secret) UpdateObject(cli kubernetes.Interface) (result *v1.Secret, err error) {
	log.Info("Updating Kubernetes Secret...")

	if k, ok := s.Data[s.ParamType]; ok {
		return nil,
		errors.New(fmt.Sprintf("Key '%s' already exists in the Secret %s/%s", k, s.Namespace, s.Name))
	}

	//s.Data[s.ParamType] = s.ParamValue
	s.Secret.StringData[s.ParamType] = s.ParamValue
	return cli.CoreV1().Secrets(s.Namespace).Update(&s.Secret)
}
