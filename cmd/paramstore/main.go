package main

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"

	//"github.com/cmattoon/aws-param-store/internal/secret"
	"github.com/cmattoon/aws-param-store/internal/param"
)



func main() {
	// Watch for annotations:
	
	//	p1 := param.NewStringParam("com.entic.bar")
	//	p2 := param.NewSecureStringParam("com.entic.foo", "")

	session, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2"),
		Credentials: credentials.NewEnvCredentials(),
	})

	svc := ssm.New(session)

	param := param.NewSecureStringParam("com.entic.foo", "")

	param.GetValue()

	
	
	// param, err := svc.GetParameter(&ssm.GetParameterInput{
// 		Name: aws.String("com.entic.foo"),
// 		WithDecryption: aws.Bool(true),
// 	})

// 	if err != nil {
// 		fmt.Println("ERROR")
// 		fmt.Println(err)
// 		return
// 	}
	
// 	fmt.Printf("%s: \"%s\"\n", *param.Parameter.Name, *param.Parameter.Value)
// 	secret_name := "my-secret"
// 	doc := heredoc.Docf(`
// ---
// apiVersion: v1
// kind: Secret
// metadata:
//   name: %s
// type: Opaque
// data:
//   %s: "%s"
// `, secret_name, *param.Parameter.Name, *param.Parameter.Value)
// 	fmt.Println(doc)
}
