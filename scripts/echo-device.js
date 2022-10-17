//
// module for echo device
//
set("123-echo-321", {
    dial9_onchange(val) {
        let echoDevice = home.getDeviceByName("echo device")
        echoDevice.set("dial9", 9)

        return home.preventUpdate
    },

    dialin_onchange(value) {
        let echoDevice = home.getDeviceByName("echo device")

        echoDevice.set("dialout",value)
        
        return home.preventUpdate
    },

})
