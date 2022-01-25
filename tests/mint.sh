#!/bin/bash

ADDRESS=${ADDRESS:-"0xd220FCc8727075B174C50eE4833CF6FE001301c6"}
HOST=${1:-"http://169.51.195.137:30000"}
NUMCALLS=${2:-"1"}
AUTHTOKEN=${3:-""}

for i in $(seq 1 ${NUMCALLS}); do
    echo "Calling ${HOST}/transaction/mint/${ADDRESS}"
	curl -X POST --header "x-auth-token: ${AUTHTOKEN}" -d @mint-payload.json ${HOST}/transaction/token/mint/${ADDRESS} &
done

wait
