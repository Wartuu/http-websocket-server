let ws

function connect() {
    var token = document.getElementById("token").value

    path = window.location.hostname
    ws = new WebSocket("ws://" + path +":7000/ws/controller", token)
}

function next() {
    ws.send(`{"action":"next", "details":"none", "conn":1}`)
}

function back() {
    ws.send(`{"action":"back", "details":"none", "conn":1}`)
}
