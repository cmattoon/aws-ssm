cmattoon/aws-ssm
================

[![Build Status](https://travis-ci.org/cmattoon/aws-ssm.svg?branch=master)](https://travis-ci.org/cmattoon/aws-ssm)

Updates Kubernetes `Secrets` with values from AWS Parameter Store


MVP Working (go binary)
-----------------------
1. Create Parameter in AWS Parameter Store

2. Create Kubernetes Secret with Annotations

3. Run Binary 

4. A key with the name `$ParameterType` should have been added to your Secret



Build
-----

    make
    make container

