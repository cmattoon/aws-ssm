package secret


import (
	b64 "encoding/base64"

	"github.com/cmattoon/aws-param-store/pkg/provider"
)

type Secret struct {
	// The Name of the secret
	Name string

	// Optional. The Ciphertext, as received from some API or whatever
	// [Ciphertext ->] Plaintext -> Value
	Ciphertext string
	
	// The Plaintext value
	// [Ciphertext ->] Plaintext -> Value
	Plaintext string
	
	// This should be a base64-encoded value for Kubernetes
	// [Ciphertext ->] Plaintext -> Value
	Value string
}

// Sets and returns Value from Plaintext
func (s Secret) GetValue() (string) {
	if s.Value == "" {
		s.Value = b64.StdEncoding.EncodeToString([]byte(s.GetPlaintext()))
	}
	return s.Value
}

// Should decrypt if needed
func (s Secret) GetPlaintext() (string) {
	return s.Plaintext
}
