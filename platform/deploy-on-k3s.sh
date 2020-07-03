#!/usr/bin/env bash

CLUSTER_NAME=parvae

k3d create -n $CLUSTER_NAME -w 3 --enable-registry
export KUBECONFIG="$(k3d get-kubeconfig --name=$CLUSTER_NAME)"

kubectl create namespace argocd
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml
