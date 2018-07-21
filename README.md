cmattoon/aws-ssm
================

[![Build Status](https://travis-ci.org/cmattoon/aws-ssm.svg?branch=master)](https://travis-ci.org/cmattoon/aws-ssm)

Updates Kubernetes `Secrets` with values from AWS Parameter Store


MVP Working (go binary)
-----------------------
1. Create Parameter in AWS Parameter Store

`my_value = foobar`

2. Create Kubernetes Secret with Annotations

```
apiVersion: v1
kind: Secret
metadata:
  name: my-secret
  annotations:
    "alpha.ssm.cmattoon.com/k8s-secret-name": my-secret
    "alpha.ssm.cmattoon.com/aws-param-name": my_value
    "alpha.ssm.cmattoon.com/aws-param-type": SecureString
    "alpha.ssm.cmattoon.com/aws-param-key": "alias/aws/ssm"
data: {}
```

3. Run Binary 

4. A key with the name `$ParameterType` should have been added to your Secret


```
apiVersion: v1
kind: Secret
metadata:
  name: my-secret
  annotations:
    "alpha.ssm.cmattoon.com/k8s-secret-name": my-secret
    "alpha.ssm.cmattoon.com/aws-param-name": my_value
    "alpha.ssm.cmattoon.com/aws-param-type": SecureString
    "alpha.ssm.cmattoon.com/aws-param-key": "alias/aws/ssm"
data:
  SecureString: foobar
```


Build
-----

    make
    make container

