# Calendar plugin



## example curl commands 
add an event
```bash
curl https://localhost:4242/plugin/calendar/addEvent -XPOST -k -d'{
    "id": "holiday",
    "created": "16-12-2022 13:19",
    "nexttrigger": "16-12-2022 17:00",
    "createdby": "me",
    "notify": ["me"],
    "msg": "holidays are comming", 
    "location": "home",
    "repeatcount": 1,
    "repeatevery": 6 
}'
```

list all known calendar events
```bash
curl https://localhost:4242/plugin/calendar/getAllEvents -XPOST -k -d'{}'
```

get event with the given event id
```bash
curl https://localhost:4242/plugin/calendar/getEvent -XPOST -k -d'{"id": "eventid"}'
```

```bash
curl https://localhost:4242/plugin/calendar/getEventByDate -XPOST -k -d'{
    "start": "2022-12-16T16:20:00Z",
    "end": "2022-12-16T17:30:00Z"
}'
```