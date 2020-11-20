#!/bin/bash
#
# WARNING: this script will wipe all applications and all namespaces managed by parvaeres
#

UUID_REGEX="[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}"

# FIXME: we should introduce an annotation instead of relying on the UUID regex
applications=$(kubectl get applications -n argocd | grep -E "$UUID_REGEX" | awk '{print $1}')
namespaces=$(kubectl get namespaces | grep -E "$UUID_REGEX" | awk '{print $1}')

DRY_RUN='--dry-run=client'
if [ "$1" == '--do-it' ] ; then
    DRY_RUN=""
fi

if [ -n "$applications" ] ; then
    echo "$applications" | xargs kubectl $DRY_RUN delete applications -n argocd
fi

if [ -n "$namespaces" ] ; then
    echo "$namespaces" | xargs kubectl $DRY_RUN delete namespaces
fi
