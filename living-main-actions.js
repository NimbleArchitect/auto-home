// module for living-main-light

function flash() {

}


// module for camera in the bar
function record_onchange(props) {
    // call the sendMsg function in a seperate thread, this is non-blocking
    // thread(home.getGroupByPath("users/home").sendUserMsg,"door opened")
    thread(sendUserMsg,"door opened")
    
    // dosent cover when one adult is home and the other is out?
    // dosent cover notification to users not at home
    if (home.getGroupByPath("users/adults").present == true) {        
        // runs when someone is home
        // set the light to red for one second
        flash_light_colour("red")
    }
}

function flash_light_colour(colour) {
    let light = home.getDeviceByName("living main")
    let previous = light.get("color")
    light.set("color",colour)
    home.sleep(1)
    light.set("state",previous)
}

function doorbell_onchange() {
    // doorbell pressed event
    let adult = home.getGroupByPath("users/adults")
    
    // need to run this as a seprate thread... how??
    adult.sendmessage("doorbell pressed at " + Date)

    // this needs to run in order
    let light = home.getDeviceByName("living main")
    light.set("state",true)
    home.sleep(1)
    light.set("state",false)
}

//
// in group module users/adults
//

////////////////////////////////////////////////////////////////////////////////////////

//
// alarm group module
//

// alarm group holds all sensor devices, any change to a sensor calls this onChange function
function group_onchange(props) {

    if (home.getDevice("alarm").get("state").asBool == true) {
        // fire alarm

        return true
    } 

    if (home.getDevice("alarm").get("downstairs").asBool == true) {
        // fire alarm

        return true
    }
}
