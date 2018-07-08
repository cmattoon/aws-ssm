package param


import (
	"log"
	"fmt"

	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/cmattoon/aws-param-store/internal/secret"
)


const (
	TypeString = "String"
	TypeStringList = "StringList"
	TypeSecureString = "SecureString"
)


type AWSParameter struct {
	Type string
	Name string
	Value string
	KmsKey string
}


func NewStringParam(name string) (AWSParameter) {
	p := AWSParameter{
		Name: name,
		Value: "",
		KmsKey: "",
		Type: TypeString,
	}
	return p
}


func NewSecureStringParam(name string, key string) (AWSParameter) {
	if key == "" {
		key = "alias/aws/ssm"
	}
	
	p := AWSParameter{
		Name: name,
		Value: "",
		KmsKey: key,
		Type: TypeSecureString,
	}
	return p
}

func NewStringListParam(name string) (AWSParameter) {
	p := AWSParameter{
		Name: name,
		Value: "",
		KmsKey: "",
		Type: TypeStringList,
	}
	return p
}


func (p Parameter) GetValue(svc ssm.Service) (string) {
	if p.Value == "" {
		param, err := svc.GetParameter(&ssm.GetParameterInput{
			Name: aws.String(p.Name),
			WithDecryption: aws.Bool((p.Type == TypeSecureString)),
		})
		
		if err != nil {
			log.Fatalf("Couldn't get value: %s", err)
		}
		p.Value = *param.Parameter.Value
	}
	return p.Value
}


func (p Parameter) ToSecret() {
	s := secret.Secret{

	}
}
