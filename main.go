package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)


type UserPool struct {
	Clients map[string]*websocket.Conn
}

type ConnectionModel struct {
    Operation string `json:"operation"`
    User string `json:"user"`
    From string `json:"from"`
    To string `json:"to"`
    Message string `json:"message"`
}


const (
    CONNECT string   = "connect"
    MESSAGE  string       = "message"
    DISCONNECT string    = "disconnect"
)

// We'll need to define an Upgrader
// this will require a Read and Write buffer size
var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
  WriteBufferSize: 1024,

  // We'll need to check the origin of our connection
  // this will allow us to make requests from our Frontend
  // development server to here.
  // For now, we'll do no checking and just allow any connection
  CheckOrigin: func(r *http.Request) bool { return true },
}

// define a reader which will listen for
// new messages being sent to our WebSocket
// endpoint
func reader(pool *UserPool, conn *websocket.Conn) {
    for {
    // read in a message
        messageType, p, err := conn.ReadMessage()
        if err != nil {
            log.Println(err)
            // deleting user when connection is down
            // it's automaticly happening from client side it's not efficient if there is a lot of user make it with goroutine # TODO
            for key, val := range pool.Clients{
                if val == conn{
                    delete(pool.Clients, key)
                    break
                }
            }
            log.Println("Active users ", pool.Clients)
            return
        }
        // print out that message for clarity
        var wsmessage ConnectionModel
        if err := json.Unmarshal(p, &wsmessage); err != nil {
            panic(err)
        }
        fmt.Println("recieved Message => ", string(p), "Message Type => ", messageType)

        switch {
            case wsmessage.Operation == CONNECT:
                log.Println(wsmessage.User, " is appending to Map of users")
                if pool.Clients == nil {
                    pool.Clients = make(map[string]*websocket.Conn)
                }
                pool.Clients[wsmessage.User] = conn
                send_msg  := map[string]string{"operation:": CONNECT, "message": "success"}
                sendMessage(conn, messageType, send_msg)
            case wsmessage.Operation == MESSAGE:
                targetUser := pool.Clients[wsmessage.To]
                sendMsg  := map[string]string{"operation:": MESSAGE, "from": wsmessage.From, "message": wsmessage.Message}
                sendMessage(targetUser, messageType, sendMsg)
            case wsmessage.Operation == DISCONNECT:
                delete(pool.Clients, wsmessage.User)
                send_msg  := map[string]string{"operation:": DISCONNECT, "message": "success"}
                sendMessage(conn, messageType, send_msg)
            default:
                log.Println("Wrong Operation Type")
        }
        log.Println("Active users ", pool.Clients)
        
       

    }
}
func sendMessage(conn *websocket.Conn,messageType int, msg map[string]string){
    newMsg, err := json.Marshal(msg)
    if err != nil {
        fmt.Println(err)
    }
    if err := conn.WriteMessage(messageType, []byte(newMsg)); err != nil {
        log.Println(err)
        return
    }
}
// define our WebSocket endpoint
func serveWs(pool *UserPool, w http.ResponseWriter, r *http.Request) {
    fmt.Println(r.Host)

    ws, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println(err)
  	}
  // listen indefinitely for new messages coming
  // through on our WebSocket connection
    reader(pool, ws) // make it with goroutine # TODO
}

func setupRoutes() {
	pool := UserPool{}

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(&pool, w, r)
	})
}


func main() {
    fmt.Println("Real Time Server v0.01")
    setupRoutes()
    http.ListenAndServe(":8088", nil)
}