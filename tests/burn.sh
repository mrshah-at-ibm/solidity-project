#!/bin/bash

TOKENID=${TOKENID:-"1"}
HOST=${1:-"169.51.195.137:30000"}
NUMCALLS=${2:-"1"}
AUTHTOKEN=${3:-""}

for i in $(seq 1 ${NUMCALLS}); do
	echo "Calling ${HOST}/transaction/token/${TOKENID}/burn"
	curl -X POST --header "x-auth-token: ${AUTHTOKEN}" ${HOST}/transaction/token/${TOKENID}/burn &
done

wait
