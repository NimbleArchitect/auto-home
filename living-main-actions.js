// module for living-main-light

function flash() {

}


// module for camera in the bar
function record_onchange(props) {
    // call the sendMsg function in a seperate thread, this is non-blocking
    thread(home.getGroupByPath("users/home").sendUserMsg,"door opened")
    
    
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
    
    const greeting = new Promise((resolve, reject) => {
        resolve("Hello!");
      });
    
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
function sendUserMsg(msg) {
    sendmessage(alert,msg)
}

function sendmessage(msg) {

    let users = home.getGroupByPath("users/adults").getUsers()
    
    for (let i = 0; i < users.length; i++) {
        if (users[i].presence == true) {
            home.plugin("telegram").sendmessage(info, msg)
        } else {
            home.plugin("telegram").sendmessage(alert, msg)
        }
    } 
}