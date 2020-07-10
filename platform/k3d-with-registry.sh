#!/usr/bin/env bash
#
# Starts a k3s cluster (via k3d) with local image registry enabled,
# and with nodes annotated such that Tilt (https://tilt.dev/) can
# auto-detect the registry.
#
# Source: https://github.com/tilt-dev/k3d-local-registry/blob/master/k3d-with-registry.sh

set -o errexit

# 🚨 only compatible with k3d v1.x (at least for now) 🚨
if ! k3d -version | grep 'v1' > /dev/null 2>&1; then
  echo "This script only works with k3d v1.x"
  exit 1
fi

# desired cluster name (default is "k3s-default")
CLUSTER_NAME="${CLUSTER_NAME:-k3s-default}"

# Check if cluster already exists.
# AFAICT there's no good way to get the registry name/port from a running
# cluster, so if it already exists, just bail.
for cluster in $(k3d ls 2>/dev/null | tail -n +4 | head -n -1 | awk '{print $2}'); do
  if [ "$cluster" == "$CLUSTER_NAME" ]; then
      # TODO(maia): check if the cluster already has the appropriate annotations--then we're okay
      # TODO(maia): if cluster exists, has registry, doesn't have annotations, apply them.
      #   (Unfortunately there's no easy way to check what registristry (if any) the cluster
      #   is running, see https://github.com/rancher/k3d/issues/193)
      echo "Cluster '$cluster' already exists, aborting script."
      echo "\t(You can delete the cluster with 'k3d delete --name=$CLUSTER_NAME' and rerun this script.)"
      exit 1
  fi
done

k3d create --enable-registry --name=${CLUSTER_NAME} "$@"

echo
echo "Waiting for Kubeconfig to be ready..."
timeout=$(($(date +%s) + 30))
until [[ $(date +%s) -gt $timeout ]]; do
  if k3d get-kubeconfig --name=${CLUSTER_NAME} > /dev/null 2>&1; then
    export KUBECONFIG="$(k3d get-kubeconfig --name=${CLUSTER_NAME})"
    DONE=true
    break
  fi
  sleep 0.2
done

if [ -z "$DONE" ]; then
  echo "Timed out trying to get Kubeconfig"
  exit 1
fi


echo
echo "Installing ArgoCD"
kubectl create namespace argocd
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml

echo
echo "Set kubecontext with:"
echo "export KUBECONFIG=\"\$(k3d get-kubeconfig --name=${CLUSTER_NAME})\""
