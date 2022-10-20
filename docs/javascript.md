
javascript is supported via the goja module as this allow me to expose the home object

# home object
the home object supports the following methods

## obj.getDevice(id)
returns the device with the specified id

### Arguments

| Parameter  | Type | Description |
| - | - | - |
| id  | String | the id of the device you are interested in |

### Returns
* device object


### Example
```javascript
home.getDevice("device-id")
```

## obj.getDeviceByName(name)
returns the device object matching the given name

### Arguments

| Parameter  | Type | Description |
| - | - | - |
| name  | String | name of the device |

### Returns
* device object


### Example
```javascript
home.getDeviceByName("device name")
```

## obj.getDeviceByPath(path)

### Arguments

| Parameter  | Type | Description |
| - | - | - |
| path  | String | path to the device in the form groupname/device |

### Returns
* device object


### Example
```javascript
home.getDeviceByPath("path/to/device")
```

## obj.getDeviceInGroup(name)
returns all devices in the group named name

### Arguments

| Parameter  | Type | Description |
| - | - | - |
| id  | String | group name |

### Returns
* device object


### Example
```javascript
home.getDeviceInGroup("group name")
```


# group object
## obj.getDevice(id)
## obj.hasDevice(id)
## obj.hasDeviceByName(id)
## obj.getDeviceByName(name)
## obj.setAll(value)
## obj.getGroup(name)
## obj.getGroupByPath(path/to/group)
## obj.getDeviceByPath(path/to/device)
## obj.getDeviceInGroup(group, device)


# device object
the device object supports the following methods

## obj.get(name)
get the property matching the specified name

### Arguments

| Parameter  | Type | Description |
| - | - | - |
| name  | String | the property name for the device you are interested in |

### Returns
* property object


### Example
```javascript
home.getDeviceByName("Kitchen Light").get("brightness")
```

## obj.set(name, value)
sets the named property to the specified value

### Arguments

| Parameter  | Type | Description |
| - | - | - |
| name  | String | the name of the property to update |
| value | any | the value to set on the named property this is device and property dependent |

### Returns
* none


### Example
```javascript
home.getDeviceByName("Kitchen Light").set("brightness", 50)

home.getDeviceByName("Kitchen Light").set("state", "on")

let echoDevice = home.getDeviceByName("echo device")
echoDevice.set("dialout",value)
```

## obj.isInGroup(name)
returns true if the device obj is in the specifiecd group

### Arguments

| Parameter  | Type | Description |
| - | - | - |
| name  | String | the name of the property to update |


### Returns
* Boolean


### Example
```javascript
home.getDeviceByName("Kitchen Light").isInGroup("kitchen")
```

# property object
the property object supports various methods and js properties based on its internal type

| js properties | Type | Description |
| - | - | - |
| value  | any | value at the time the event was fired |
| latest | any | the current "live" value, this is pulled at the time it is called not before so expect a delay on heavey systems |
| previous | any | the value before this event was triggered |

## obj.last(x)
searches the current property history and returns the item number

### Arguments

| Parameter  | Type | Description |
| - | - | - |
| x  | Integer | the history item number to retrieve |

### Returns
* any - a single value dependent on the internal property type


### Example
```javascript
home.getDeviceByName("Kitchen Light").get("brightness").last(3)

let echoDevice = home.getDeviceByName("echo device")
echoDevice.get("dialout").last(2)
```

## obj.asPercent()

### Returns
* Integer

### Example
```javascript
```

## obj.isSwitch()

### Returns
* Boolean

### Example
```javascript
```

## obj.isDial()

### Returns
* Boolean

### Example
```javascript
```

## obj.isButton()

### Returns
* Boolean

### Example
```javascript
```

## obj.isText()

### Returns
* Boolean

### Example
```javascript
```
