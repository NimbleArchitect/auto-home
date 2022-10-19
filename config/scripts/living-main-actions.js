//
// module for living-room-light
//
set("123-tv-light-321", {
    // only called when brightness actually changes from its 
    // current value
    hue_onchange(val) {
        console.log("hue change is at your service :)")
        // console.log(">>"+ val)

        let light = home.getDeviceByName("living main")

        light.set("state", true)
        home.sleep(1)
        light.set("state", false)

    },

    state_onchange(value) {
        console.log("state change is at your service :)")
        console.log("!>>" + value)

        let light = home.getDeviceByName("living main")

        light.set("state", true)
        home.sleep(1)
        light.set("state", false)

        console.log("1!>" + home.getDeviceByName("TV light").get("hue").value)
        console.log("2!>" + home.getDeviceByName("TV light").get("hue").latest)
        console.log("3!>" + home.getDeviceByName("TV light").get("hue").previous)

    },

    flash() {

    },

    set_colour() {

    }
})

set("group/kitchen-sockets", {
    onchange(val) {
        console.log("group called kitchen-sockets")
    }
})

set("group/living-room-sockets", {
    onchange(val) {
        console.log("group called living-room-sockets")
    }
})

set("group/living-room", {
    onchange(val) {
        console.log("group called living-room")
    }
})
