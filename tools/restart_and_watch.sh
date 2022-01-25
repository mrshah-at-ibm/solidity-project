#!/bin/bash

#cd ..
#./build-quick.sh
#cd -
kubectl delete po -n mrshah --all
sleep 3
kubectl logs -f -n mrshah -l app=mrshah-app
