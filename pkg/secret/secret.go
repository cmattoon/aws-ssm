package secret

type GeneratedSecret struct {
	AWS *AWSParameter
	K8S *KubernetesSecret
}

type SecretSpec struct {

}

func NewGeneratedSecret(spec *SecretSpec) GeneratedSecret {
	return GeneratedSecret{}
}
