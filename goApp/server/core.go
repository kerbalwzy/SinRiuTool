package server

import (
	"flag"
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":41001", "GoApp后台服务地址")
var upGrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logger.Info(r.RemoteAddr, r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "www/html/index.html")
}

func WsRequest(w http.ResponseWriter, r *http.Request) {
	conn, err := upGrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error:%v", err)
		return
	}
	client := SingleNewClient(conn)
	logger.Info("new client connect", client.conn.RemoteAddr())
}

func StartLocalServer() {
	flag.Parse()
	router := httprouter.New()
	router.GET("/", Index)
	router.ServeFiles("/static/*filepath", http.Dir("www"))
	router.GET("/ws", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		WsRequest(writer, request)
	})
	logger.Info("start local sever:", *addr)
	err := http.ListenAndServe(*addr, router)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
