cmattoon/aws-ssm
================

[![Build Status](https://travis-ci.org/cmattoon/aws-ssm.svg?branch=master)](https://travis-ci.org/cmattoon/aws-ssm)

Updates Kubernetes `Secrets` with values from AWS Parameter Store

 * For example usage, see `example.yaml`
 * Use the Helm chart to get up and running quickly


Helm Chart
----------

Use `AWS_REGION=<region> ./install_chart.sh` to install from source


### Install Script Environment

The following environment variables must be set for `install_chart.sh`:

  - `AWS_REGION`
  - `AWS_ACCESS_KEY`
  - `AWS_SECRET_KEY`


### Values
| Value        | Default          | Example                     | Description                                                      |
|--------------|------------------|-----------------------------|------------------------------------------------------------------|
| aws_region   |                  | us-west-2                   | The AWS region in which the Pod is deployed                      |
| kubeconfig64 |                  | <string>                    | The output of `$(cat $KUBE_CONFIG | base64)`. Stored as a Secret |
| metrics_port | 9999             | <int>                       | Serve metrics/healthchecks on this port                          |
| replicas     | 1                | <int>                       | The number of Pods                                               |
| image.name   | cmattoon/aws-ssm | <docker-repo>/<image-name>  | The Docker image to use for the Pod container                    |
| image.tag    | latest           | <docker-tag>                | The Docker tag for the image                                     |
| resources    | {}               | <dict>                      | Kubernetes Resource Requests/Limits                              |
|              |                  |                             |                                                                  |


Docker Container
----------------

### Build

Run `make container` to build the Docker image


Configuration
-------------

The following values can be provided via environment variables or CLI flags.
CLI flags take precdence over environment variables

| Environment | Flag         | Default        | Description                      |
|-------------|--------------|----------------|----------------------------------|
| KUBE_CONFIG | -kube-config |                | The path to the kube config file |
| MASTER_URL  | -master-url  |                | The Kubernetes master API URL    |
| METRICS_URL | -metrics-url | localhost:9999 | Address for healthchecks/metrics |
| AWS_REGION  | -region      | us-west-2      | The AWS Region                   |


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

