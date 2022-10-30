//
// common functions
//

function sendUserMsg(msg) {
    sendmessage(alert,msg)
}

function sendmessage(msg) {
    console.log("called sendmessage >>")
    plugin.telegram.sendMessage(msg)
    
}
