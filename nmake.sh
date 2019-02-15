#!/usr/bin/env bash
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
sudo chown -R willard:willard ${DIR}/static
pushd ${DIR}
make image
docker tag containous/traefik gitlab-registry.nordstrom.com/gtm/linkerd-sandbox/traefik:v1.7$1
docker push gitlab-registry.nordstrom.com/gtm/linkerd-sandbox/traefik:v1.7$1
popd

