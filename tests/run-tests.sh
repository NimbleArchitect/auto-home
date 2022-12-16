#!/bin/bash

cd tests

############################
out=$(go run ./request-response-timing/*.go 2> /dev/null)
if [[ $? -ne 0 ]]; then
    echo "$out"
    exit 1
fi
echo "$out" | grep longest

############################
out=$(go run ./request-response-echo/*.go 2> /dev/null)
if [[ $? -ne 0 ]]; then
    echo "$out"
    exit 1
fi
echo "$out"

############################
out=$(go run ./high-event-count/*.go 2> /dev/null)
if [[ $? -ne 0 ]]; then
    echo "$out"
    exit 1
fi
echo "$out"





cd ..