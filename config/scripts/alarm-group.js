//
// alarm group module
//
set("group/alarm", {
    // alarm group holds all sensor devices, any change to a sensor calls this onChange function
    onchange(props) {
        console.log("alarm group triggered")

        thread("sendmessage","++++++++++TESTY MESSAGE++++++++++")

        if (home.getDeviceByName("alarm").get("state").asBool() == true) {
            // fire alarm
            console.log("FIRE1 FIRE2 FIRE3!!!")

            return home.stopProcessing
        } else {
            console.log(">> not yet!!!")
        }

        if (home.getDeviceByName("alarm").get("downstairs").asBool() == true) {
            // fire alarm
            console.log("FIRE FIRE FIRE!!!")

            return home.stopProcessing
        }
    }
})
