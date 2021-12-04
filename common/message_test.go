package common

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestMessageA(t *testing.T) {

	m := NewJSONMessage(1, 2, []byte("Hello World"))
	fmt.Println(string(m.Bytes()))
	rm := NewJSONMessageByBytes(m.Bytes())
	fmt.Println(rm.From(), rm.To(), rm.Timestamp(), string(rm.Data()))

}

func TestMessageB(t *testing.T) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	fmt.Println("开始运行...")
	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		if ws, err := upgrader.Upgrade(rw, r, nil); err == nil {
			connection := NewWebSocketConnection(ws)
			defer connection.Close()

			for {
				if data, err := connection.Read(); err == nil {
					if recvM := NewJSONMessageByBytes(data); recvM != nil {
						t := &time.Time{}
						t.Add(time.Duration(recvM.Timestamp()) * time.Second)
						datetime := t.Format("2006-01-02 15:04:05")
						fmt.Println(recvM.From(), recvM.To(), string(recvM.Data()), datetime)
						connection.Write(NewJSONMessage(recvM.To(), recvM.From(), []byte("RECVIVED IT!")).Bytes())
					}
				} else {
					fmt.Println(err)
					break
				}
			}

		} else {
			fmt.Println(err)
		}
	})
	http.ListenAndServe(":8080", nil)

}
