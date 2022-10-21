# ** currently in development, don't expect much **
interested and want to help, please drop a message in [discussions](https://github.com/NimbleArchitect/auto-home/discussions)

# Introduction
A fully reactive home automation controller, this is probably the third revision of an idea that started back in 2018

designed to work on a raspberry pi and be network aware the new code uses JavaScript as the internal scripting language which allows you to link up onchange and ontrigger events, its also possible to create custom functions that control devices, i.e. if you want to flash your lights red you can by creating a JavaScript function that flashes the lights for you. 

The advanced group support allows devices to be part of multiple groups and groups to also be part of other groups so you can end up like this ascii image, I have labelled the groups and devices to make it a bit easier to read.

```
├── kitchen (group)
│ ├── kitchen fan (device)
│ ├── kitchen lights (group)
│ │ ├── main kitchen light (device)
│ │ ├── under cupboard light (device)
│ │ └── island light (device)
│ └── kitchen sockets (group)
│ ├── socket 1 (device)
│ ├── socket 2 (device)
│ └── island socket (device)
└── downstairs lights (group)
├── kitchen lights (group)
└── downstairs lights (group)
```

So... On its own and having just the devices grouped is pretty pointless until we start attaching events to the groups, lets say we want the kitchen fan to be turned on when ever a light is also turned on. Now in most smart home systems you would have to set an event for every light in the kitchen and if you have lots of lights or more complicated requirements it can be easy to make a mistake. Due to the way this system bubbles events up to groups this becomes as simple as setting an onchange event on the kitchen lights group that will turn the fan on and off, this then leaves the device onchange event free to set the brightness of the light depending on the time of day.
Its also completly possible to script the kitchen lights group to automatically adjust the brightness of every light within its group based on the time of day or other events, in short its a powerful system with endless possibilities that works for your lifestyle.

# Features
* fully reactive system
* multi threaded and multi processor capable
* JavaScript scripting engine, scripts can be attached to devices groups and users
* plugin system allows for custom devices
* multiple processing pipelines allow for events to be processed simultaneously
* device events also trigger their owning group
* custom functions allow you to write your own device responses/actions
* virtual devices with event support
* event repeat protection, devices/groups support a cool down period of X milliseconds where duplicate events are ignored
* server start script
* scripts can prevent group events and device updates from happening

# Planned features
* web UI
* simplified event programmer
* per screen UI
* ability to record and replay events
* support for file uploads
* full http API support
* user presence support
* more plugins (sunrise/sunset, )
* integrated calendar
* custom plugin interfaces
* call http post/get from JavaScript
* devices support custom data fields that can be read/written from JavaScript
* custom fields are automatically removed after a set age
* plus much, much more...

# Building and running

to build change into the auto-home folder and go build and run with ./server

# Further information
I have a rough outline of how the system works in the [design](./docs/design.md) documentation

for specific information on how to write the script files consult the [scripts](./scripts/) folder and the [javascript](./docs/javascript.md) document. 
General JavaScript programming examples can be found all over the internet

# Contributions, features and improvements
if you want to help drop a message in [discussions](https://github.com/NimbleArchitect/auto-home/discussions)

raise a ticket for feature requests and improvements

# License
TBD

