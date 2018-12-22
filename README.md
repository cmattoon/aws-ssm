cmattoon/aws-ssm
================

[![Build Status](https://travis-ci.org/cmattoon/aws-ssm.svg?branch=master)](https://travis-ci.org/cmattoon/aws-ssm)
![GitHub tag](https://img.shields.io/github/tag/cmattoon/aws-ssm.svg)
![Docker Pulls](https://img.shields.io/docker/pulls/cmattoon/aws-ssm.svg)
[![codecov](https://codecov.io/gh/cmattoon/aws-ssm/branch/master/graph/badge.svg)](https://codecov.io/gh/cmattoon/aws-ssm)
[![Go Report Card](https://goreportcard.com/badge/github.com/cmattoon/aws-ssm)](https://goreportcard.com/report/github.com/cmattoon/aws-ssm)
[![Maintainability](https://api.codeclimate.com/v1/badges/764dddb334f5dc9fb986/maintainability)](https://codeclimate.com/github/cmattoon/aws-ssm/maintainability)
[![Anchore Image Overview](https://anchore.io/service/badges/image/7d144c4a4e096c3f87c563080ea1279aed19e718ccdf12a6b7436e086090d3b3)](https://anchore.io/image/dockerhub/cmattoon%2Faws-ssm%3Alatest)


Updates Kubernetes `Secrets` with values from AWS Parameter Store

 * For example usage, see `example.yaml`
 * Use the Helm chart to get up and running quickly

Build Options
-------------

  * Helm Chart (recommended): `make {lint|install|purge}`
  * Go: `make test && make build`
  * Docker: `make container`


Helm Chart
----------

### Install Helm Chart

First, export required variables, then run `make install`.


    export AWS_REGION=<region>
    export AWS_SECRET_KEY=<secret>
    export AWS_ACCESS_KEY=<access-key-id>


### AWS User/Role

The AWS credentials should be associated with an IAM user/role that has the following permissions:

  - @todo
  

### Values

The following chart values may be set. Only the required variables (AWS credentials) need provided by the user. Most of the time, the other
defaults should work as-is.


| Req'd | Value          | Default          | Example                     | Description                                                      |
|-------|----------------|------------------|-----------------------------|------------------------------------------------------------------|
| YES   | aws.region     | ""               | us-west-2                   | The AWS region in which the Pod is deployed                      |
| YES   | aws.access_key | ""               |                             |                                                                  |
| YES   | aws.secret_key | ""               |                             |                                                                  |
| NO    | kubeconfig64   | ""               | <string>                    | The output of `$(cat $KUBE_CONFIG | base64)`. Stored as a Secret |
| NO    | metrics_port   | 9999             | <int>                       | Serve metrics/healthchecks on this port                          |
| NO    | replicas       | 1                | <int>                       | The number of Pods                                               |
| NO    | image.name     | cmattoon/aws-ssm | <docker-repo>/<image-name>  | The Docker image to use for the Pod container                    |
| NO    | image.tag      | latest           | <docker-tag>                | The Docker tag for the image                                     |
| NO    | resources      | {}               | <dict>                      | Kubernetes Resource Requests/Limits                              |
| NO    | host_ssl_dir   | ""               | /etc/ssl/certs              | If specified, mounts certs from the host.                        |
| NO    | rbac.enabled   | true             | <bool>                      | Whether or not to add Kubernetes RBAC stuff                      |


Docker Container
----------------

### Build

Run `make container` to build the Docker image


Configuration
-------------

The following app config values can be provided via environment variables or CLI flags.
CLI flags take precdence over environment variables.

A KUBE_CONFIG and MASTER_URL are only necessary when running outside of the cluster (e.g., dev)

| Environment | Flag         | Default        | Description                      |
|-------------|--------------|----------------|----------------------------------|
| AWS_REGION  | -region      | us-west-2      | The AWS Region                   |
| METRICS_URL | -metrics-url | 0.0.0.0:9999   | Address for healthchecks/metrics | 
| KUBE_CONFIG | -kube-config |                | The path to the kube config file |
| MASTER_URL  | -master-url  |                | The Kubernetes master API URL    |


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


```
apiVersion: v1
kind: Secret
metadata:
  name: my-secret
  annotations:
    "alpha.ssm.cmattoon.com/k8s-secret-name": app-secrets
    "alpha.ssm.cmattoon.com/aws-param-name": /path/to/env
    "alpha.ssm.cmattoon.com/aws-param-type": Directory
    "alpha.ssm.cmattoon.com/aws-param-key": "alias/aws/ssm"
data:
  file1: value1
  file2: value2
```


Build
-----

    make
    make container


CA Certificates
---------------

For ease of use, the `ca-certificates` package is installed on the final `library/alpine` image. If you're having SSL/TLS
connection issues, `export HOST_SSL_DIR=/etc/ssl/certs` before running `make install`. This will mount the SSL cert directory
on the EC2 instance.