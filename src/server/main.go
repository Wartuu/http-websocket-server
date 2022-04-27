package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/websocket"
	"server/data"
)

type ControllerMsg struct {
	Action  string `json:"action"`
	Details string `json:"details"`
	Conn    int    `json:"conn"`
}

type ClientMsg struct {
	Action  string `json:"action"`
	Details string `json:"details"`
}

type ClientResponse struct {
	Success bool   `json:"success"`
	Info    string `json:"info"`
	IsAlive int    `json:"alive"`
}

type pingResp struct {
	Alive int `json:"alive"`
}

var PORT string = ":7000"
var ToClientInput string

var BroadcastMessage string
var BroadcastID int64 = 0

func Controller(writer http.ResponseWriter, request *http.Request) {
	log.Println(request.Header.Get("Sec-Websocket-Protocol"))
	if request.Header.Get("Sec-Websocket-Protocol") == "7000::" {
		log.Println("New Controller. Sending socket")
		ws := websocket.Server{Handler: websocket.Handler(ControllerWs)}
		ws.ServeHTTP(writer, request)
	} else {
		log.Println("Wrong auth")
	}
}

func Client(writer http.ResponseWriter, request *http.Request) {
	log.Println("New Client. Sending socket")
	ws := websocket.Server{Handler: websocket.Handler(ClientWS)}
	ws.ServeHTTP(writer, request)
}

func ClientWS(ws *websocket.Conn) {
	Connections++
	log.Println("new Connection succed!!!!")
	go isAlive(ws)

	OLD := BroadcastID

	for {
		time.Sleep(1000 * time.Millisecond)
		if OLD != BroadcastID {
			log.Println("SENDING TO CLIENT!")
			websocket.Message.Send(ws, BroadcastMessage)
			OLD = BroadcastID
			break
		}

	}
}

func ControllerWs(ws *websocket.Conn) {
	log.Println("NEW CONTROLLER")
	for {
		var privateMsg ControllerMsg
		var message string
		websocket.Message.Receive(ws, &message)
		json.Unmarshal([]byte(message), &privateMsg)

		if privateMsg.Conn == 0 {
			ws.Close()
			return
		}
		if privateMsg.Conn != 0 {
			if privateMsg.Action == "back" {
				num--
			} else if privateMsg.Action == "next" {
				num++
			}
			log.Println(message)
			numstr := strconv.Itoa(num)
			BroadcastMessage = strings.Join([]string{"{\"Action\":\"switchnum\", \"Details\":\"", numstr, "\"}"}, "")
			BroadcastID++

		}
	}
}

func isAlive(ws *websocket.Conn) {
	var ResponsePing ClientResponse
	var ResponseString string
	for {
		time.Sleep(10 * time.Second)
		websocket.Message.Send(ws, "{\"Action\":\"ping\",\"Details\":\"1\"}")

		websocket.Message.Receive(ws, &ResponseString)
		json.Unmarshal([]byte(ResponseString), &ResponsePing)
		if ResponsePing.IsAlive == 0 {
			Connections--
			ws.Close()
			break
		}

	}
}

func ControllerSite(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(writer, string(data.Controller))
}

func ClientSite(writer http.ResponseWriter, request *http.Request) {
	switch request.URL.Query().Get("page") {
	case "":
		fmt.Fprintf(writer, "%s" ,string(data.P1))
	case "1":
		fmt.Fprintf(writer, "%s" ,string(data.P1))
	case "2":
		fmt.Fprintf(writer, "%s" ,string(data.P2))
	case "3":
		fmt.Fprintf(writer, "%s" ,string(data.P3))
	case "4":
		fmt.Fprintf(writer, "%s" ,string(data.P4))
	case "5":
		fmt.Fprintf(writer, "%s" ,string(data.P5))
	case "6":
		fmt.Fprintf(writer, "%s" ,string(data.P6))
	case "7":
		fmt.Fprintf(writer, "%s" ,string(data.P7))
	case "8":
		fmt.Fprintf(writer, "%s" ,string(data.P8))
	case "9":
		fmt.Fprintf(writer, "%s" ,string(data.P9))
	case "10":
		fmt.Fprintf(writer, "%s" ,string(data.P10))
	default:
		fmt.Fprintf(writer, "%s" ,string(data.P1))
	}

}

func ControllerJS(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "text/javascript; charset=utf-8")
	fmt.Fprintf(writer, "%s", string(data.ControllerJS))
}

func ClientCSS(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "text/css; charset=utf-8")
	fmt.Fprintf(writer, "%s", string(data.ClientCSS))
}

func ClientJS(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "text/javascript; charset=utf-8")
	fmt.Fprintf(writer, "%s" ,string(data.ClientJS))
}

func main() {
	go ConnectionList()
	http.HandleFunc("/ws/client", Client)
	http.HandleFunc("/", ClientSite)
	http.HandleFunc("/ws/controller", Controller)
	http.HandleFunc("/controller", ControllerSite)

	http.HandleFunc("/js/controller", ControllerJS)
	http.HandleFunc("/js/client", ClientJS)

	http.HandleFunc("/css/client", ClientCSS)

	log.Println("Running at 127.0.0.1" + PORT)

	err := http.ListenAndServe(PORT, nil)

	if err != nil {
		panic(err.Error())
	}
}

var Connections int64 = 0

func ConnectionList() {
	for {
		time.Sleep(500 * time.Millisecond)
		fmt.Println(Connections)
	}

}

var num = 0
