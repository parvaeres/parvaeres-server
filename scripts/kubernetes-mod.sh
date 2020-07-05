#!/usr/bin/env bash
#
# Workaround issues importing argocd library importing k8s stuff
# See: https://github.com/kubernetes/kubernetes/issues/79384#issuecomment-521493597
#
#
# Here's an example of this issue:
#
#  ➜ go test ./pkg/gitops 
# go: finding module for package github.com/argoproj/argo-cd/pkg/apis/application/v1alpha1
# go: finding module for package github.com/smartystreets/goconvey/convey
# go: finding module for package github.com/google/uuid
# go: found github.com/argoproj/argo-cd/pkg/apis/application/v1alpha1 in github.com/argoproj/argo-cd v1.6.1
# go: found github.com/google/uuid in github.com/google/uuid v1.1.1
# go: found github.com/smartystreets/goconvey/convey in github.com/smartystreets/goconvey v1.6.4
# go: github.com/argoproj/argo-cd@v1.6.1 requires
#    k8s.io/kubernetes@v1.16.6 requires
#    k8s.io/api@v0.0.0: reading k8s.io/api/go.mod at revision v0.0.0: unknown revision v0.0.0
#

set -euo pipefail

# version should match argocd
# See: https://github.com/argoproj/argo-cd/blob/master/go.mod
VERSION=${1#"v"}
if [ -z "$VERSION" ]; then
    echo "Must specify version!"
    exit 1
fi
MODS=($(
    curl -sS https://raw.githubusercontent.com/kubernetes/kubernetes/v${VERSION}/go.mod |
    sed -n 's|.*k8s.io/\(.*\) => ./staging/src/k8s.io/.*|k8s.io/\1|p'
))
for MOD in "${MODS[@]}"; do
    V=$(
        go mod download -json "${MOD}@kubernetes-${VERSION}" |
        sed -n 's|.*"Version": "\(.*\)".*|\1|p'
    )
    go mod edit "-replace=${MOD}=${MOD}@${V}"
done
go get "k8s.io/kubernetes@v${VERSION}"
