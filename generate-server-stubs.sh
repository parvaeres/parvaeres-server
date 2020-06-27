#!/usr/bin/env bash
#
# Generate server stubs for ParvaeRes API
# 

docker run \
    --rm \
    --user "$(id -u):$(id -g)" \
    --volume "${PWD}:/local" \
    --workdir "/local" \
    --env GIT_USER_ID=ContainerSolutions \
    --env GIT_REPO_ID=parvaeres \
    openapitools/openapi-generator-cli \
        generate \
        -i parvaeres-api.yaml \
        -c openapi-config.yaml \
        -g go-experimental \
        -o parvaeres-server
