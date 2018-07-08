package secret


import (
	b64 "encoding/base64"
	"github.com/MakeNowJust/heredoc"
)


type Secret struct {
	
	ParamName string
	ParamValue string
	
	SecretName string
	SecretKey string
	SecretValue string
	
}

func NewSecret() Secret {
	s := Secret{}
	return s
}

func FromParam(p *param.AWSParameter) {
	s := NewSecret()
	s.ParamName
}

func (s Secret) GetManifest() (string) {
	return heredoc.Docf(`
---
apiVersion: v1
kind: Secret
metadata:
  name: %s
type: Opaque
data:
  %s: "%s"
`, s.SecretName, s.SecretKey, s.GetEncodedValue())	
}

func (s Secret) GetEncodedValue() (string) {
	if s.ParamValue == "" {
		return ""
	}
	return b64.StdEncoding.EncodeToString([]byte(s.ParamValue))
}
