# connections
first you must provide a user and pass to /v1/connect the server will return a session id that can be used to access various endpoints

## request
```bash
curl -k -XPOST https://127.0.0.1:4242/v1/connect -d '{"data":{"user":"user.test.user","pass":"secretpassword"}}'
```
## response
```json
{
    "result": {
        "status": "ok",
        "msg": ""
    },
    "data": {
        "session": "11111111-1111-1111-1111-111111111111"
    }
}
```


# api calls
the following paths all require a valid session id, replace $sessionid with the session id returned from /v1/connect

## /v1/device
returns a list of all known devices along with their properties and current state
```bash
curl -k -H "session: $sessionid" https://127.0.0.1:4242/v1/device
```

## /v1/device/deviceid
returns device information and all properties along with their known state
```bash
curl -k -H "session: $sessionid" https://127.0.0.1:4242/v1/device/123-echo-321
```

## /v1/device/deviceid/propertyid
returns the known property state from the device with id deviceid
```bash
curl -k -H "session: $sessionid" https://127.0.0.1:4242/v1/device/123-echo-321/dialdelay
```

## /v1/device/deviceid/propertyid?setstate=value
sets the property of device to value
```bash
curl -k -H "session: $sessionid" https://127.0.0.1:4242/v1/device/123-echo-321/dialdelay?setstate=181
```

## /v1/device/deviceid/propertyid?setpercent=percentage
sets the property of device to percentage
```bash
curl -k -H "session: $sessionid" https://127.0.0.1:4242/v1/device/123-echo-321/dialdelay?setpercent=25
```
