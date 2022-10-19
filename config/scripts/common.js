//
// common functions
//

function sendUserMsg(msg) {
    sendmessage(alert,msg)
}

function sendmessage(msg) {
    console.log("called sendmessage >>")
    home.plugin("Telegram").call("SendMessage", {
        message: msg, 
    })
    
}
