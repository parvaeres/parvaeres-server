#!/bin/bash
#
# NOTE: run this script from root of the project repository checkout.
#

K3D_VERSION_TAG=v3.0.1
K3D_INSTALL_DIR=.
K3D_INSTALL_URL=https://raw.githubusercontent.com/rancher/k3d/main/install.sh
K3D_CLUSTER_NAME=${K3D_CLUSTER_NAME:-parvaeres}

# make sure we find the k3d executable if we already downloaded it
export PATH=$PATH:$K3D_INSTALL_DIR/


# download k3d if not available
download_k3d() {
    if ! command -v k3d &> /dev/null ; then
        curl -s $K3D_INSTALL_URL | USE_SUDO=false \
                                   TAG=$K3D_VERSION_TAG \
                                   K3D_INSTALL_DIR=$K3D_INSTALL_DIR bash
    fi
}

# cluster and registry up!
up() {
    # create the cluster if doesn't exist
    if ! k3d cluster list "$K3D_CLUSTER_NAME" &> /dev/null ; then
        k3d cluster create "$K3D_CLUSTER_NAME" \
            --agents 3 \
            --wait \
            --timeout 120s \
            --update-default-kubeconfig \
            --switch-context \
            --volume "$(pwd)/platform/k3s/registries.yaml:/etc/rancher/k3s/registries.yaml"
    fi

    # start a local docker registry if not running
    if ! docker ps --format '{{.Names}}' | grep -w 'registry\.localhost' &> /dev/null; then
        docker volume create local_registry
        docker container run -d --rm \
            --name registry.localhost \
            -v local_registry:/var/lib/registry \
            -p 5000:5000 registry:2

        # connect the registry to the cluster's network
        docker network connect "k3d-$K3D_CLUSTER_NAME" registry.localhost
    fi
}

# cluster and registry down!
down() {
    download_k3d

    if docker ps --format '{{.Names}}' | grep -w 'registry\.localhost' &> /dev/null; then
        docker container rm -f registry.localhost
        docker volume rm -f local_registry
    fi

    if k3d cluster list "$K3D_CLUSTER_NAME" &> /dev/null ; then
        k3d cluster delete "$K3D_CLUSTER_NAME"
    fi

    if k3d network list "k3d-$K3D_CLUSTER_NAME" &> /dev/null ; then
        k3d network rm "k3d-$K3D_CLUSTER_NAME"
    fi
}

case $1 in
up)
    download_k3d
    up
    ;;
down)
    download_k3d
    down
    ;;
*)
    echo "$0 up | down"
    ;;
esac
