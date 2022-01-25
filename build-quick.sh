#!/bin/bash -e

./scripts/build.sh

docker build -t mrshah2/kp:app -f Dockerfile-quick .

docker push mrshah2/kp:app
