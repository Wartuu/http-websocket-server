path = window.location.hostname
console.log(path)
var ws = new WebSocket("ws://" + path +":7000/ws/client")

var contentdisplay = document.getElementById("presentation")


ws.onmessage = (data) => {
    console.log(data.data)
    parsedJSON = JSON.parse(data.data)

    switch (parsedJSON.Action) {
        case "ping":
            ws.send(`{"success":"true", "info":"none", "alive":1}`)
            break
        case "switchnum":
            window.location.replace("?page="+ parsedJSON.Details)
            break
    }

}