package main

import (
	"encoding/json"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/websocket"
	"github.com/kataras/neffos"
	"log"
	"strconv"
)

type App struct {
	app *iris.Application
	ws  *neffos.Server
}

func main() {
	application := App{
		app: iris.New(),
	}
	app := application.app

	app.HandleDir("/", "./www")

	app.Get("/test", func(c context.Context) {
		_, _ = c.WriteString("hello iris")
	})

	app.Get("/send", func(c context.Context) {
		send(application.ws)
		_, _ = c.WriteString("send ok")
	})

	application.ws = socket()

	id := 0
	app.Get("/echo", websocket.Handler(application.ws, func(c context.Context) string {
		id += 1
		return strconv.Itoa(id)
	}))

	err := app.Listen(":8081")
	if err != nil {
		log.Fatal("Listen", err)
	}
}

func socket() *neffos.Server {
	ws := websocket.New(websocket.DefaultGobwasUpgrader, websocket.Events{
		websocket.OnNativeMessage: func(nsConn *websocket.NSConn, msg websocket.Message) error {
			log.Printf("Server got: %s from [%s]", msg.Body, nsConn.Conn.ID())
			return nil
		},
	})
	ws.OnConnect = func(c *neffos.Conn) error {
		log.Println("新的连接...", c.ID())
		return nil
	}
	ws.OnDisconnect = func(c *neffos.Conn) {
		log.Println(c.ID(), "断开连接...")
	}

	return ws
}

func send(ws *neffos.Server) {
	data := iris.Map{
		"a": "1",
		"b": "2",
		"c": "3",
	}
	bytes, _ := json.Marshal(data)
	for _, connection := range ws.GetConnections() {
		_ = connection.Socket().WriteText(bytes, 0)
		_ = connection.Socket().WriteBinary(bytes, 0)
		//_ = connection.Socket().WriteText([]byte(connection.ID()+" hello world"), 0)
		//connection.Write(neffos.Message{
		//	Body: []byte(connection.ID() + " hello world"),
		//})
	}
}
