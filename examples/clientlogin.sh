#!/bin/bash

# login to server an retrieve session id
connect_response=$(curl -k -XPOST https://127.0.0.1:4242/v1/connect -d '{"data":{"user":"user.test.user","pass":"secretpassword"}}')

# check result status
status=$(echo $connect_response |jq -r '.result.status')
if [[ $status == "ok" ]]; then
    # result status was ok so we can grab the session id
    sessionid=$(echo $connect_response |jq -r '.data.session')

    # use the session id for the curl command
    # to get all known devices and their state
    curl -k -H "session: $sessionid" https://127.0.0.1:4242/v1/device

    # get all peoperties of device 123-echo-321
    # curl -k -H "session: $sessionid" https://127.0.0.1:4242/v1/device/123-echo-321

    # get dialdelay peoperty of device 123-echo-321
    # curl -k -H "session: $sessionid" https://127.0.0.1:4242/v1/device/123-echo-321/dialdelay

fi
