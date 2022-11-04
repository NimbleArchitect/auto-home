
JavaScript is implemented via the goja module as this allow me to expose the home object

# overview
all scripts are loaded during server start from the scripts folder (pointed to via config.json), scripts are loaded as an encapsulated object to prevent rouge code polluting the runtime environment. 

## set(id, obj)
Currently the set command is exposed to allow registration of objects into the device scope, to use pass in the id and an object

### Arguments
| Parameter  | Type | Description |
| - | - | - |
| id  | String | the id of the device you want to attach the supplied object to the following are also acceptable "group/id", "device/id", "user/id" where id is replaced with the id you want to register against|
| obj | Object | the object to attach, the objects methods are exposed |


### Returns
* none

### Example
to register a device use something like the following but remember to replace device-id with the id of you device and property with the name of the property you want to trigger against.

```javascript
// don’t forget to change device-id and property below
set("device-id", {
    // ontrigger is optional
    property_ontrigger(value) {
        // code goes here
    },

    // onchange is also optional
    property_onchange(value) {
        // code here
    },
})

// or you can also use device/device-id
// dont forget to replace device-id
set("device/device-id", {
    // ontrigger is optional 
    property_ontrigger(value) {
        // code goes here
    },

    // onchange is also optional
    property_onchange(value) {
        // code here
    },
})
```
lets say you have a device called "kitchen light" with an id of "123_kitchen-light" you would replace "device-id" with "123_kitchen-light" and if the light has a property called brightness and you want to run a script when the brightness changes you would replace "property_onchange(value)" with "brightness_onchange(value)" so you end up with the following

```javascript
set("123_kitchen-light", {
    brightness_onchange(value) {
        // code goes here, value holds the new value that
        //  was received by the server
    },
})
```
you can also register onchange events with groups by prefixing "group/" in front of the group name, here we have a group called "kitchen" that holds a device called "kitchen light" and by using the below we can register for "onchange" events of the named group
```javascript
set("group/kitchen", {
    onchange() {
        // code goes here
    },
})
```
now using the above when any property of the kitchen light changes this group is also called with its attached onchange event the same works for any device or group that is a member of the kitchen group 

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
| path  | String | path to the device in the form groupid/device |

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

## obj.sleep(seconds)
pauses the script for the specified number of seconds

### Arguments

| Parameter  | Type | Description |
| - | - | - |
| seconds  | Integer | number of seconds to pause the script |

### Returns
* none

### Example
```javascript
home.sleep(5)
```

## obj.countdown(name, milliseconds, function)
creates or restarts a timer called name, name is a timer identifier so must be unique, when the number of miliseconds have elapsed the timer calls function, if countdown is called again before the timer reaches zero the timer is reset.

to remove the timer set milliseconds to 0

### Arguments

| Parameter  | Type | Description |
| - | - | - |
| name  | String | unique countdown name |
| milliseconds  | Integer | length of time until the countdown timer calls function |
| function  | Function | valid javascript function to run |

### Returns
* none

### Example
```javascript
// after 60 seconds turn the light off
home.countdown("porchlight", 60000, function () {
    light.set("state", false) // turn off the light
})
```
```javascript
// to disable set the second parameter to 0
home.countdown("porchlight", 0)
```
```javascript
// while the porch door is open keep the light on, once the door is closed turn the light off after 60 seconds
// if the door is opened again before 60 seconds have passed we keep the light on
light.set("state",true) // first we turn the light on
 while (porchdoor.get("state").latest == "open") { // while door is open
    if (light.get("state").asBool() == true) { // and the light is on
        // set the countdown timer, caling this in a loop means we reset before 60 seconds have passed
        home.countdown("porchlight", 60000, function () {
            // when the timer expires
            light.set("state", false) // turn off the light
        })

        home.sleep(5) // this stops the system running in a tight loop and abusing resources
    }
}

```

# group object
all group methods are called the same way as they are from the home object 

| js properties | Type | Description |
| - | - | - |
| name  | string | the name of the selected group |
| id | string | the current group id |


## obj.getDevice(id)
> see home.getDevice()

## obj.hasDevice(id)
> see home.hasDevice()

## obj.hasDeviceByName(id)
> see home.hasDeviceByName()

## obj.getDeviceByName(name)
> see home.getDeviceByName()

## obj.setAll(value)
> see home.setAll()

## obj.getGroup(name)
> see home.getGroup()

## obj.getGroupByPath(path/to/group)
> see home.getGroupByPath()

## obj.getDeviceByPath(path/to/device)
> see home.getDeviceByPath()

## obj.getDeviceInGroup(group, device)
> see home.getDeviceInGroup()


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
| latest | any | the current "live" value, this is pulled at the time it is called not before so expect a delay on heavy systems |
| previous | any | the value before this event was triggered |

## obj.last(x)
searches the current property history and returns item x from the internal array

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
returns the stored value as a percentage

### Returns
* Integer - between 0 and 100

### Example
```javascript
home.getDeviceByName("Kitchen Light").get("brightness").asPercent()
```

## obj.isSwitch()
returns true if the property is an internal switch type, false otherwise

### Returns
* Boolean

### Example
```javascript
home.getDeviceByName("Kitchen Light").get("brightness").isSwitch()
```

## obj.isDial()
returns true if the property is an internal dial type, false otherwise

### Returns
* Boolean

### Example
```javascript
home.getDeviceByName("Kitchen Light").get("brightness").isDial()
```

## obj.isButton()
returns true if the property is an internal button type, false otherwise

### Returns
* Boolean

### Example
```javascript
home.getDeviceByName("Kitchen Light").get("brightness").isButton()
```

## obj.isText()
returns true if the property is an internal text type, false otherwise

### Returns
* Boolean

### Example
```javascript
home.getDeviceByName("Kitchen Light").get("brightness").isText()
```