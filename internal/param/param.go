package param

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

