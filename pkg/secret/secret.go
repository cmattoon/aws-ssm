package secret

import (
	log "github.com/sirupsen/logrus"
	"github.com/cmattoon/aws-ssm/pkg/provider"
)

type Secret struct {
	Name string
	Namespace string
	Values map[string]string
}


func NewSecret(p provider.Provider, name string, ns string, decrypt bool) (*Secret) {
	s := &Secret{
		Name: name,
		Namespace: ns,
		Values: map[string]string{},
	}
	value, err := p.GetParameterValue(name, decrypt)
	if err != nil {
		log.Infof("Couldn't get value for %s: %s", name, err)
	} else {
		s.Values["value"] = value
	}
	return s
}
