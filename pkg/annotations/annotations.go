/**
 * Copyright 2018 Curtis Mattoon
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package annotations

const (
	K8SSecretName = "alpha.ssm.cmattoon.com/k8s-secret-name"
	K8SSecretType = "alpha.ssm.cmattoon.com/k8s-secret-type"
	AWSParamName  = "alpha.ssm.cmattoon.com/aws-param-name"
	AWSParamType  = "alpha.ssm.cmattoon.com/aws-param-type"
	AWSParamKey   = "alpha.ssm.cmattoon.com/aws-param-key"
	AWSRoleArn    = "iam.amazonaws.com/role"

	V1ParamName = "aws-ssm/aws-param-name"
	V1ParamType = "aws-ssm/aws-param-type"
	V1ParamKey  = "aws-ssm/aws-param-key"
	V1RoleArn   = "aws-iam/role-arn"
)
