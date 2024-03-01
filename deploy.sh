#!/usr/bin/env bash

set -e

script_dir="$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"

pnpm run build

rsync -avr --exclude node_modules --exclude .git "${script_dir}/" admin@192.168.0.115:/apps/notes/

kubectl apply -f ./deployment.template.yaml

kubectl -n notes rollout restart deployment deployment
