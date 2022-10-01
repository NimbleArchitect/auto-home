// module for living-room-light

// only called when brightness actually changes from its 
// current value
function hue_onchange(val) {
    console.log("hue change is at your service :)")
    // console.log(">>"+ val)

    let light = home.getDeviceByName("living main")

    light.set("state",true)
    home.sleep(1)
    light.set("state",false)
    
}

//called everytime brightness is triggered
// function hue_ontrigger(val) {
//     console.log("brightness trigger is at your service :)")
//     console.log(">>"+ val)
// }

//called during sllent updates
function brightness_onupdate() {
    
}

// function state_ontrigger(val) {
//     console.log("state trigger is at your service :)")
//     console.log(">>"+ val)

// }

function state_onchange(value) {
    console.log("state change is at your service :)")
    console.log("!>>"+ value)

    let light = home.getDeviceByName("living main")
    // for (let i = 0; i < 254; i=i+20) {
    //     light.set("brightness", i)
    //     home.sleep(1)
    // } 

    // home.sleep(2)
    
    // for (let i = 254; i >0; i=i-20) {
    //     light.set("brightness", i)
    //     home.sleep(1)
    // } 

    light.set("state",true)
    home.sleep(1)
    light.set("state",false)
    

    // home.getDeviceByName("TV light").set("state", "off")

    console.log("1!>" + home.getDeviceByName("TV light").get("hue").value)
    console.log("2!>" + home.getDeviceByName("TV light").get("hue").latest)
    console.log("3!>" + home.getDeviceByName("TV light").get("hue").previous)

    //
    // console.log("sleep completed")
}

function flash() {

}

function set_colour() {

}


