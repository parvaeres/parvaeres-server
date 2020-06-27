#!/usr/bin/env bash
#
# Generate server stubs for ParvaeRes API
# 

# docker run -ti --rm openapitools/openapi-generator-cli $@ ; exit 0

DESTINATION=parvaeres-server

docker run \
    --rm \
    --user "$(id -u):$(id -g)" \
    --volume "${PWD}:/local" \
    --workdir "/local" \
    openapitools/openapi-generator-cli \
        generate \
        --input-spec parvaeres-api.yaml \
        --config openapi-config.yaml \
        --generator-name go-server \
        --git-user-id riccardomc \
        --git-repo-id parvaeres \
        --output $DESTINATION

gofmt -w $DESTINATION
goimports -w $DESTINATION
