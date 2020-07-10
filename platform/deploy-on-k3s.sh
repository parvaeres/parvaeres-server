#!/usr/bin/env bash

CLUSTER_NAME=parvae

k3d create -n $CLUSTER_NAME -w 3 --enable-registry
export KUBECONFIG="$(k3d get-kubeconfig --name=$CLUSTER_NAME)"

