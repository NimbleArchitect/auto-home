#!/bin/bash

cd tests

############################
out=$(./request-response-timing/test-run 2> /dev/null)
if [[ $? -ne 0 ]]; then
    echo "$out"
    exit 1
fi
echo "$out" | grep longest

############################
out=$(./request-response-echo/test-run 2> /dev/null)
if [[ $? -ne 0 ]]; then
    echo "$out"
    exit 1
fi
echo "$out"




cd ..