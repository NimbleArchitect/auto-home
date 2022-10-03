//
// alarm group module
//

// alarm group holds all sensor devices, any change to a sensor calls this onChange function
function group_onchange(props) {

    if (home.getDevice("alarm").get("state").asBool == true) {
        // fire alarm

        return home.StopProcessing
    } 

    if (home.getDevice("alarm").get("downstairs").asBool == true) {
        // fire alarm

        return home.StopProcessing
    }
}
