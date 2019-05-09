cmattoon/aws-ssm
================

[![Build Status](https://travis-ci.org/cmattoon/aws-ssm.svg?branch=master)](https://travis-ci.org/cmattoon/aws-ssm)
![GitHub tag](https://img.shields.io/github/tag/cmattoon/aws-ssm.svg)
![Docker Pulls](https://img.shields.io/docker/pulls/cmattoon/aws-ssm.svg)
[![codecov](https://codecov.io/gh/cmattoon/aws-ssm/branch/master/graph/badge.svg)](https://codecov.io/gh/cmattoon/aws-ssm)
[![Go Report Card](https://goreportcard.com/badge/github.com/cmattoon/aws-ssm)](https://goreportcard.com/report/github.com/cmattoon/aws-ssm)
[![Maintainability](https://api.codeclimate.com/v1/badges/764dddb334f5dc9fb986/maintainability)](https://codeclimate.com/github/cmattoon/aws-ssm/maintainability)


Updates Kubernetes `Secrets` with values from AWS Parameter Store

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


### AWS Credentials

Uses the [default credential provider chain](https://docs.aws.amazon.com/sdk-for-go/api/aws/credentials/#NewChainCredentials)


### Values

The following chart values may be set. Only the required variables (AWS credentials) need provided by the user. Most of the time, the other
defaults should work as-is.


| Req'd | Value          | Default          | Example                     | Description                                                      |
|-------|----------------|------------------|-----------------------------|------------------------------------------------------------------|
| YES   | aws.region     | ""               | us-west-2                   | The AWS region in which the Pod is deployed                      |
| NO    | aws.access_key | ""               |                             | REQUIRED when no other auth method available (e.g., IAM role)    |
| NO    | aws.secret_key | ""               |                             | REQUIRED when no other auth method available (e.g., IAM role)    |
| NO    | kubeconfig64   | ""               | <string>                    | The output of `$(cat $KUBE_CONFIG \| base64)`. Stored as a Secret|
| NO    | metrics_port   | 9999             | <int>                       | Serve metrics/healthchecks on this port                          |
| NO    | image.name     | cmattoon/aws-ssm | <docker-repo>/<image-name>  | The Docker image to use for the Pod container                    |
| NO    | image.tag      | latest           | <docker-tag>                | The Docker tag for the image                                     |
| NO    | resources      | {}               | <dict>                      | Kubernetes Resource Requests/Limits                              |
| NO    | rbac.enabled   | true             | <bool>                      | Whether or not to add Kubernetes RBAC stuff                      |
| NO    | ssl.mount_host | false            | <bool>                      | Mounts {ssl.host_path} -> {ssl.mount_path} as hostVolume         |
| NO    | ssl.host_path  | /etc/ssl/certs   | <path>                      | The SSL certs dir on the host                                    |
| NO    | ssl.mount_path | /etc/ssl/certs   | <path>                      | The SSL certs dir in the container (dev)                         |


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
| LOG_LEVEL   | -log-level   | info           | The Logrus log level             |


Basic Usage
-----------
1. Create Parameter in AWS Parameter Store

`my-db-password` = `foobar`

2. Create Kubernetes Secret with Annotations

```
apiVersion: v1
kind: Secret
metadata:
  name: my-secret
  annotations:
    aws-ssm/k8s-secret-name: my-secret
    aws-ssm/aws-param-name: my-db-password
    aws-ssm/aws-param-type: SecureString
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
    aws-ssm/k8s-secret-name: my-secret
    aws-ssm/aws-param-name: my-db-password
    aws-ssm/aws-param-type: SecureString
data:
  SecureString: Zm9vYmFyCg==
```

Annotations
-----------

| Annotation                 | Description                                            | Default         |
|----------------------------|--------------------------------------------------------|-----------------|
| `aws-ssm/k8s-secret-name`  | The name of the Kubernetes Secret to modify.           | `<none>`        |
| `aws-ssm/aws-param-name`   | The name of the AWS SSM Parameter. May be a path.      | `<none>`        |
| `aws-ssm/aws-param-type`   | Determines how values are parsed, if at all.           | `String`        |
| `aws-ssm/aws-param-key`    | Required if `aws-ssm/aws-param-type` is `SecureString` | `alias/aws/ssm` |


### AWS Parameter Types

Values for `aws-ssm/aws-param-type` are:

| Value          | Behavior                 | AWS Value                   | K8S Value(s)                            |
|----------------|--------------------------|-----------------------------|-----------------------------------------|
| `String`       | No parsing is performed  | `foo` = `bar`               | `foo: bar`                              |
| `SecureString` | Requires `aws-param-key` | `foo` = `bar`               | `foo: bar`                              |
| `StringList`   | Splits CSV mapping       | `foo=bar,bar=baz,baz=bat`   | `foo: bar`<br> `bar: baz`<br>`baz: bat` |
| `Directory`    | Get multiple values      | `/path/to/values`           | <treats each subkey/value as a String>  |



Build
-----

    make           # Build binary
    make container # Build Docker image
    make push      # Push Docker image


CA Certificates
---------------

For ease of use, the `ca-certificates` package is installed on the final `library/alpine` image. If you're having SSL/TLS
connection issues, `export HOST_SSL_DIR=/etc/ssl/certs` before running `make install`. This will mount the SSL cert directory
on the EC2 instance.