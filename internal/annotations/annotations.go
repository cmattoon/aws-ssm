package annotations

const (
	// The name of the param in SSM
	AWSParameterName = "alpha.ssm.cmattoon.com/param-name"
	
	// String, StringList, SecureString
	AWSParameterType = "alpha.ssm.cmattoon.com/param-type"
	
	// Key as in crypto key
	AWSParameterKey = "alpha.ssm.cmattoon.com/param-key"

	// The 'metadata.name' of the Kubernetes Secret
	K8SecretName = "alpha.ssm.cmattoon.com/secret-name"
	
	// Key as in "key-value"
	K8SecretKey = "alpha.ssm.cmattoon.com/secret-key"
)
