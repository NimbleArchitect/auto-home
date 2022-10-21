# design

at its core auto-home is based on two loops with and attached buffer, as its written in go I took advantage of channels which allows the system to be fully responsive in reacting to incoming events.

We might as well start our walkthrough at the point the event enters the server

## http3
I wanted something that could communicate across the network without using TCP for a long time I toyed a few times with UDP but could never get it work how I wanted.

Until I came across the execllent quic-go package which the server now uses for its communication as it allows for multiple streams of data in a bi directional fashion all while being performant with poor connections.

for now I wont talk about the clients except to say that an event message is sent from the client code using a http3 post request to /event/registrationid the body of which is sent as json. The server recieves the event message does some basic validation and sends the message to the eventManager calling the AddEvent function

## eventManager
When AddEvent is called we add the message to the chAdd channel this starts our first part of the processing loop the EventManager. The event manager recieves the message and attempts to save the message in the ring buffer (events array) once successfull the eventCount is increased and the array id is sent to the chCurrentEvent channel for processing by the EventLoop. Once recieved on the channel a new go routine is spun up and the homeManager Trigger function is called with the message id, timestamp and a list of device properties as arguments only when the trigger function returns is the array id released from the EventLoop and sent back via the chRemove channel, at which point the EventManager recieves the array id and decreases the eventCount, thus completing the work of the eventManager.

## homeManager
During server startup the homeManager loads all *.js files in the script folder and compiles them to a compiledScript array, the server then starts the selected number of javascript VM's and runs all the compiled scripts on each VM, once successful the VMs are referenced in an activeVMs array and wait to be called.

On recieving the Trigger from the eventManager the homeManager performs the following tasks:

* get the next avaliable javascript virtual machine
* copies the current state of all known devices from the deviceManager
* make the plugins avaliable to javascript
* prepares the event message for the javascript VM

now the javascript VM is ready we call its Process function, once Process completes we record the event in the history logs and release the javascript VM back in to the pool

## javascript VM
### Process
At this point the javascript Process has a copy of all devices, a list of plugins and a list of the device properties that will change. So we setup an empty device object with empty properties and start the first part of our 3 step process.

Step 1 is to validate all properties passed in and check if the device properties have a named ontrigger script we can call.

for example: if the property has a name brightness then the brightness_ontrigger function is searched for within the device scripts and if found the script is run.

Step 2 we get a list of groups that the device is a member of and recursivly look for all parent groups then, we look for and call the onchange event for every group in the list, when the groups functions have been called the return code is checked to see if we should continue processing further groups or continue our merry way.

Step 3 (assuming we go this far) is our onchange process, this part of the process allows us to loop over each property that was set to be changed, first we check to make sure the stop processing flag has not been set then we check for and run the onchange script for every property before we writing the chage back to the deviceManager via the liveDevice field.

### devices
A snapshot copy of all devices is taken everytime the javascript VM is called as this allows uninterrupted access to every device in the system, it also allows script users to view the state at the time the event was triggered,and we dont have to deal with mutex locking or race conditions, I also provide access to the live (most recent) device property which can be retrieved via latest

### home
the home object is designed to be a one stop object for access to every part of the system, currently it holds a list of all known devices, group and plugins along with a mini history for every device property, it also exposes various functions to allow interaction with these items, furtner details can be found in docs/javascript.md 

## deviceManager
the device manager holds a list of all currently registered devices with a list of all properties, all properties store a mutex value for locking and a repeat window  aling with a value histoy

### repeat window
when the device generates an event it provides a list of properties that have changed each property in the list has an expire time caluclated from the currnet time plus a duration, this expire time is stored as the repeatWindowTimeStamp, then when the property is triggered again the timestamp is compared against the current time, if the expire time is more recent than the current time the property is skipped over for processing, this allow the first property change event to be processed whilst ignoring the others within a set repeat window. 

### value history
a simple ring buffer that copies the current value into the next space on the ring, this max buffer size can be set on a per device basis

## device properties
early on in the development process I decided that it was possible to express device state with a list of properties, these properties have become the building blocks of the system and are detailed below.
devices can hold multiple properties of each type with the only requirement being that they all have different names.

### switch
a simple switch, internally is stored as a simple true/false but can be set using on/off, up/down, open/closed, yes/no and true/false
the switch stores the string used and will return both the string value and its internal boolean value

### dial
stores a current value, along with a minimum and maximum value, if the value is set outside of min and max it is internall reset back to min and max respectivly 

### button
a simple push button its actually very similer to a switch except that never actually changes state events are triggered still but the state is never changed from its initial value.

### text
a free use text field can be used to store any needed required text value

## server startup


## plugins

## clients
