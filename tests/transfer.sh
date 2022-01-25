#!/bin/bash

TOKENID=${TOKENID:-"2"}
HOST=${1:-"http://169.51.195.137:30000"}
NUMCALLS=${2:="1"}
AUTHTOKEN=${3:-""}

for i in $(seq 1 ${NUMCALLS}); do
	echo "Calling ${HOST}/transaction/token/${TOKENID}/transfer"
	curl -X POST --header "x-auth-token: ${AUTHTOKEN}" -d @transfer-payload.json ${HOST}/transaction/token/${TOKENID}/transfer &
done

wait
