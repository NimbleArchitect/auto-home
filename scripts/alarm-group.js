//
// alarm group module
//
set("group/alarm", {
    // alarm group holds all sensor devices, any change to a sensor calls this onChange function
    onchange(props) {
        console.log("alarm group triggered")
        // console.log("2>>"+home.getDeviceByName("alarm").get("state").latest)
        // console.log("3>>"+home.getDeviceByName("alarm").get("state").value)
        console.log("1>>")
        console.log("2>>"+home.getDeviceByName("alarm").get("state").asBool())
        console.log("3>>")
        console.log("4>>" + (home.getDeviceByName("alarm").get("state").asBool() == true))
        console.log("5>>")

        thread(sendmessage,"++++++++++TESTY MESSAGE++++++++++")
        // home.sleep(5)
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
