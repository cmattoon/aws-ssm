#!/bin/bash
# Required values
AWS_REGION=${AWS_REGION:?"Must specify AWS_REGION"}

# Optional Values
# (probably don't need changed)
CHART_DIR=aws-ssm
RELEASE_NAME=${RELEASE_NAME:-"aws-ssm"}
RELEASE_NAMESPACE=${RELEASE_NAMESPACE:-"kube-system"}
KUBE_CONFIG=${KUBE_CONFIG:-~/.kube/config}
EXTRA_ARGS=${EXTRA_ARGS:-""}
# Get base64-encoded kube config
if [ ! -f "$KUBE_CONFIG" ]; then
    >&2 echo " [!] Missing KUBE_CONFIG ($KUBE_CONFIG)"
    exit 1;
fi
KUBECONFIG64=$(cat $KUBE_CONFIG | base64)

helm upgrade --install $RELEASE_NAME \
     --namespace $RELEASE_NAMESPACE \
     --set aws_region=$AWS_REGION \
     --set kubeconfig64="$KUBECONFIG64" \
     $EXTRA_ARGS \
     $CHART_DIR
