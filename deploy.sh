#!/usr/bin/env bash

set -e

scp ./nginx.prod.conf will@192.168.0.120:/mnt/storage/k3s-apps/notes/

docker build -t 192.168.0.120:5000/notes:latest .
docker push 192.168.0.120:5000/notes:latest
kubectl apply -f ./deployment.template.yaml
kubectl rollout restart deployment notes
